Imaging Service — Usage

Ringkasan

Imaging Service adalah service HTTP sederhana untuk mengambil gambar dari URL, menerapkan beberapa transformasi (resize, flip, filter, kualitas/format), dan mengembalikan gambar yang sudah diproses.

Kontrak singkat

- Endpoint: GET /
- Path format: /OPTIONS/ENCODED_IMAGE_URL
  - OPTIONS: string berisi instruksi transformasi (lihat bagian `Options` di bawah)
  - ENCODED_IMAGE_URL: URL gambar yang sudah di-URL-encode atau path tanpa skema (http/https akan ditambahkan secara default)
- Response: image (Content-Type sesuai format: image/jpeg, image/png, atau image/webp) atau HTTP error

Environment

- PORT (opsional): port untuk menjalankan server (default: 8080)
- File `go.mod` sudah mencantumkan dependensi yang diperlukan (imaging, webp)

Build & Run

Di Windows PowerShell pada root project:

```powershell
# build
go build ./...

# run (PORT optional)
$env:PORT = "8080"; .\imaging-service.exe
# atau
go run main.go
```

Format OPTIONS

Beberapa contoh opsi yang didukung (diletakkan sebelum slash yang memisahkan URL):

- Size: WIDTHxHEIGHT
  - Contoh: 200x150
  - Jika width < 0, maka gambar akan dibalik horizontal (flip). Contoh: -200x150
- Smart crop: tambahkan teks `smart` di opsi untuk menandakan smart crop (catatan: fitur smart crop belum diimplementasikan di processor)
- Crop region: X:Y:W:H (empat angka, dipisah dengan `:`)
  - Contoh: 10:20:300:200
- Filters: tambahkan `filters:` diikuti daftar filter yang dipisah `:` dengan format name(value)
  - Contoh: filters:grayscale(1):blur(2.5):brightness(10):contrast(5):format(webp):quality(80):watermark(MyWatermark)
  - Filter yang dikenali oleh parser/processor saat ini: grayscale, blur, brightness, contrast
  - `format(...)` dan `quality(...)` di-parse dari filter string untuk menentukan output format dan kualitas

Contoh URL

1) Resize menjadi 300x200 dan ambil gambar dari example.com

http://localhost:8080/300x200/https%3A%2F%2Fexample.com%2Fimage.jpg

2) Resize + blur + output WEBP kualitas 80

http://localhost:8080/300x200/filters:blur(2.5):format(webp):quality(80)/https%3A%2F%2Fexample.com%2Fimage.jpg

3) Flip horizontal (negatif width)

http://localhost:8080/-300x200/https%3A%2F%2Fexample.com%2Fimage.jpg

Catatan implementasi & batasan saat ini (hal-hal yang belum bekerja / risiko)

1) Crop region parsing ada, tetapi cropping belum diterapkan di `processor.ProcessImage`.
   - Lokasi: `internal/parser/optionParser.go` (mengisi CropRegion), `internal/processor/imageProcessor.go` (tidak ada pemrosesan crop).
   - Dampak: opsi crop tidak berpengaruh.
   - Saran perbaikan: implementasikan cropping dengan `imaging.Crop` menggunakan nilai `opts.CropRegion`.

2) Smart crop (opsi `smart`) ter-deteksi di parser tetapi tidak diimplementasikan di processor.
   - Lokasi: parser menandai `opts.SmartCrop`, tapi processor tidak memeriksa.
   - Saran: gunakan strategi sederhana (center crop) atau integrasikan library deteksi fokus untuk smart crop.

3) [✓] Watermark sudah diimplementasikan
   - Cara penggunaan: tambahkan filter watermark di URL, contoh: filters:watermark(Copyright 2025)
   - Menggunakan font Go Regular dengan ukuran relatif terhadap lebar gambar
   - Watermark ditampilkan di pojok kanan bawah dengan warna putih semi-transparan

4) Parser `filters:` berbahaya bila format `( ... )` tidak ditemukan — akses `valStr[1]` digunakan tanpa pengecekan panjang sehingga dapat menyebabkan panic saat parsing jika format unexpected.
   - Lokasi: `internal/parser/optionParser.go` bagian parsing filters.
   - Dampak: malformed filter string dapat menyebabkan panic atau nilai tak terduga.
   - Saran: tambahkan pemeriksaan `len(valStr) == 2` sebelum mengakses `valStr[1]` dan fallback yang aman.

5) Banyak operasi jaringan / IO tidak memiliki timeout atau retry.
   - `pkg/utils/fetcher.go` menggunakan `http.Get` tanpa timeout.
   - Dampak: request yang tergantung bisa menggantung goroutine server.
   - Saran: gunakan `http.Client{Timeout: ...}` dan batasan ukuran body.

6) Tidak ada caching.
   - Setiap permintaan mengambil gambar dari origin; ini bisa menyebabkan latensi tinggi dan beban traffic.
   - Saran: tambahkan cache (in-memory LRU atau cache disk) untuk URL yang sering diminta.

7) Tidak ada logging request/metrics.
   - Saran: tambah logging (request path, duration, status) dan metrics (Prometheus) bila diperlukan.

8) Tidak ada validasi atau sanitasi lengkap pada URL input.
   - Parser menambahkan `https://` jika skema tidak ada — ini bisa menyembunyikan input path yang salah.
   - Saran: lebih ketat memvalidasi host/URL atau optional allowlist.

9) Tidak ada test unit / CI.
   - Saran: tambahkan beberapa unit test untuk parser dan processor.

10) Pemrosesan format `webp` bergantung pada library `github.com/chai2010/webp` yang sudah ada di `go.mod`. Pastikan environment build mendukung library tersebut (tidak perlu CGO secara default untuk library ini).

Langkah perbaikan prioritas (rekomendasi)

1. Perbaiki parsing filters untuk mencegah index out of range.
2. Implementasikan crop dan watermark di `processor.ProcessImage`.
3. Tambahkan timeout pada FetchImage (gunakan http.Client dengan Timeout).
4. Tambahkan unit tests untuk `parser.ParseOptions` dan `processor.ProcessImage`.
5. Tambahkan caching sederhana jika diperlukan untuk performa.

Penutup

Saya sudah men-parse kode sumber dan membuat ringkasan ini. Jika mau, saya bisa:
- Mengimplementasikan perbaikan prioritas (contoh: aman-guard parsing + cropping) dan menambahkan unit test singkat.
- Menambahkan example Dockerfile atau compose untuk deployment.

---

(Di-generate otomatis berdasarkan isi file project saat ini)
