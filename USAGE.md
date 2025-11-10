## Imaging Service — Usage

Ringkasan

Imaging Service adalah layanan HTTP kecil yang mengambil sebuah gambar dari URL, menerapkan transformasi (resize, crop, flip, filter, watermark), lalu mengembalikan gambar hasilnya. Dokumentasi ini menjelaskan format path, opsi yang didukung, cara build/run, contoh penggunaan, serta catatan penting.

## Endpoint

- Method: GET
- Path: /{OPTIONS}/{ENCODED_IMAGE_URL}
  - OPTIONS: string instruksi transformasi (lihat bagian "Opsi" di bawah)
  - ENCODED_IMAGE_URL: URL sumber yang sudah di-URL-encode (mis. hasil dari `encodeURIComponent`), atau host/path tanpa skema (server akan menambahkan `https://` jika skema tidak ada)

Contoh ringkas: GET /300x200/https%3A%2F%2Fexample.com%2Fimg.jpg

Response: body berisi gambar yang sudah diproses. Content-Type ditentukan berdasarkan opsi format (image/jpeg, image/png, image/webp).

## Environment

- PORT (opsional): port untuk menjalankan server. Default: 8080

## Build & Run (Windows PowerShell)

Jalankan dari root project:

```powershell
# build semua paket
go build ./...

# run (opsional set PORT)
$env:PORT = "8080"; go run main.go
```

Catatan: program menuliskan beberapa log debug ke stdout (parser/processor). Untuk produksi, sebaiknya ganti dengan logger berlevel.

## Opsi (FORMAT dari {OPTIONS})

OPTIONS adalah string yang menggabungkan beberapa instruksi, dipisahkan dengan `:`. Struktur yang sering dipakai:

- Resize: WIDTHxHEIGHT
  - Contoh: `300x200` — hasil ukuran 300x200
  - Jika WIDTH < 0, gambar akan di-flip horizontal dan width dianggap absolut. Contoh: `-300x200` = flip horizontal + resize
- Crop region: X:Y:W:H
  - Contoh: `10:20:310:220` — memotong rectangle yang dimulai di (10,20) sampai (310,220)
  - Implementation: parser mengisi `CropRegion` dan processor akan melakukan `imaging.Crop` jika region valid
- Smart crop: menyertakan kata `smart` di OPTIONS menandai permintaan smart-crop (parser mengenali `smart`, namun implementasi smart crop belum canggih — saat ini belum ada deteksi fokus)
- Filters: gunakan prefix `filters:` diikuti list filter yang dipisah `:`. Setiap filter dapat memiliki parameter dalam tanda kurung `name(value)`.
  - Contoh lengkap: `filters:blur(2.5):grayscale(1):format(webp):quality(80):watermark(https://example.com/wm.png)`
  - Filter yang di-parse dan diterapkan saat ini: `grayscale`, `blur`, `brightness`, `contrast`.
  - `format(...)` menentukan output (`jpeg`/`png`/`webp`). Default output adalah `jpeg`.
  - `quality(...)` menentukan kualitas untuk encoder JPEG/WEBP (default parser: 75).
  - `watermark(...)` dapat berisi URL (atau teks tergantung cara penggunaan); saat ini kode men-download image watermark jika `opts.Watermark` berisi URL.

Contoh OPTIONS gabungan: `300x200:filters:blur(2.5):format(webp):quality(80)`

## Contoh lengkap

1) Resize 300x200 dari example.com

GET http://localhost:8080/300x200/https%3A%2F%2Fexample.com%2Fimage.jpg

2) Resize + blur + output webp kualitas 80

GET http://localhost:8080/300x200:filters:blur(2.5):format(webp):quality(80)/https%3A%2F%2Fexample.com%2Fimage.jpg

3) Flip horizontal (negatif width)

GET http://localhost:8080/-300x200/https%3A%2F%2Fexample.com%2Fimage.jpg

4) Menambahkan watermark (contoh URL watermark)

GET http://localhost:8080/300x200:filters:watermark(https%3A%2F%2Fexample.com%2Fwm.png)/https%3A%2F%2Fexample.com%2Fimage.jpg

## Catatan implementasi & batasan (yang saya temukan)

1) Parser berpotensi panic saat parsing filter yang malformed
   - Penyebab: parsing mencari grup `(...)` dan mengambil index tanpa pemeriksaan panjang hasil regex. Jika filter tidak mengandung `(...)` seperti diharapkan, akses index bisa menyebabkan out-of-range.
   - Dampak: permintaan dengan filter malformat dapat menyebabkan panic di server.
   - Rekomendasi: tambahkan pemeriksaan hasil regex (len == 2) sebelum mengakses dan fallback ke nilai default.

2) FetchImage menggunakan `http.Get` tanpa timeout
   - Dampak: request ke origin yang lambat dapat menggantung goroutine dan menurunkan throughput.
   - Rekomendasi: gunakan `http.Client{Timeout: time.Second * 10}` dan batasi ukuran body (mis. baca sampai N bytes) serta handling kode status != 200.

3) Crop sudah diimplementasikan di `processor.ProcessImage` (periksa `CropRegion`) — kode akan melakukan crop jika region valid.

4) Smart crop di-deteksi oleh parser (`opts.SmartCrop`), namun tidak ada strategi smart-crop yang nyata di processor — saat ini tidak berpengaruh selain flag.

5) Watermark: kode saat ini mengunduh watermark dari URL bila `opts.Watermark` berisi URL. Watermark di-resize ke 1/5 lebar gambar lalu di-overlay di sudut kanan-bawah dengan alpha 0.5.

6) Tidak ada caching, logging per-request, atau metrics. Untuk penggunaan produksi, pertimbangkan menambahkan cache (LRU), logger terstruktur, dan metrics (Prometheus).

7) Tidak ada unit tests di repo saat ini. Direkomendasikan menambahkan test untuk `parser.ParseOptions` dan `processor.ProcessImage`.

## Troubleshooting singkat

- Jika server tidak mau jalan: pastikan module dependencies terunduh (`go mod tidy`) dan tidak ada error build.
- Jika image tidak muncul: periksa URL sumber dan pastikan dapat diakses dari mesin yang menjalankan service. Periksa log debug untuk pesan dari parser/processor.
- Pastikan gcc sudah terpasang di server anda, ketik gcc --version