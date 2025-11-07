Panduan Penggunaan — Imaging Service

Layanan HTTP ini berfungsi untuk mengambil gambar dari sumber eksternal, menerapkan transformasi seperti pengubahan ukuran (resize), penerapan filter, atau pembalikan arah (flip), kemudian mengembalikannya sebagai gambar hasil pemrosesan.

Menjalankan Server

Dari direktori utama proyek, jalankan perintah berikut:

# Mengunduh atau memperbarui dependensi
go mod tidy

# Menjalankan secara langsung
go run .

# Atau membangun executable, kemudian menjalankan
go build .
.\imaging-service.exe


Secara bawaan, server akan berjalan pada port 8080.

Format Endpoint

Semua permintaan dilakukan menggunakan metode GET pada root path (/).
Path digunakan untuk menyandikan (encode) parameter pemrosesan dan URL gambar.

Format umum:

/[OPTIONS]/[IMAGE_URL]


Segmen pertama (OPTIONS) berisi instruksi pemrosesan gambar.

Segmen kedua (IMAGE_URL) merupakan alamat gambar sumber (dapat di-URL-encode apabila mengandung karakter khusus).

Contoh

http://localhost:8080/400x300:filters:grayscale/picsum.photos/id/237/200/300

http://localhost:8080/400x300:filters:blur(5)/picsum.photos/id/237/200/300

Catatan:
Server tidak mendukung format dengan parameter query seperti ?url=....
Apabila path kosong, server akan mengembalikan status 400 (Bad Request).

Opsi yang Didukung
1. Pengubahan Ukuran (Resize)

Gunakan format WIDTHxHEIGHT pada segmen pertama path.
Contoh: /400x300/

400x300 → mengubah ukuran menjadi 400×300 piksel.

-400x300 → mengubah ukuran dan membalik gambar secara horizontal.

400x0 → menyesuaikan lebar saja.

0x400 → menyesuaikan tinggi saja.

2. Pemotongan Gambar (Crop) dan Smart Crop

Parser mengenali:

Token smart untuk mengaktifkan pemotongan otomatis (SmartCrop).

Format manual x1:y1:x2:y2 untuk menentukan koordinat pemotongan, disimpan dalam Options.CropRegion.

Fitur ini telah dikenali oleh parser, namun belum diimplementasikan pada tahap pemrosesan gambar.

3. Filter

Filter ditambahkan menggunakan format:

filters:NAME(VALUE):NAME2(VALUE2)


Contoh:

filters:blur(10):grayscale


Nilai dalam tanda kurung bersifat opsional (default 1.0 apabila tidak diisi).

Filter yang Telah Didukung
Nama Filter	Nilai	Deskripsi
grayscale	–	Mengubah gambar menjadi hitam putih
blur(value)	Float	Menambahkan efek kabur (semakin besar nilai, semakin kabur)
brightness(value)	Float	Menyesuaikan tingkat kecerahan
contrast(value)	Float	Menyesuaikan tingkat kontras
Filter yang Sudah Dikenali Parser namun Belum Diimplementasikan

watermark(text/url)

format(type)

quality(value)

4. Kualitas (Quality)

Nilai kualitas gambar keluaran secara bawaan adalah 75.
Nilai ini dapat diubah melalui:

filters:quality(85)


Nilai tersebut akan digunakan ketika gambar dikodekan ke format JPEG.

Contoh Penggunaan

Mengubah ukuran, membalik horizontal, dan menambahkan efek blur:

http://localhost:8080/-400x300:filters:blur%285%29/picsum.photos/id/237/200/300


Mengubah ukuran dan menambahkan dua filter (blur dan grayscale):

http://localhost:8080/400x300:filters:blur(2):grayscale/picsum.photos/id/237/200/300


Apabila URL gambar tidak mencantumkan protokol (http:// atau https://), parser akan secara otomatis menambahkan https://.
    
Respons

Body: Berisi gambar hasil pemrosesan dalam bentuk biner.

Header Content-Type: Saat ini, server selalu mengembalikan gambar dengan format image/jpeg.

Encoding: Menggunakan nilai quality yang diperoleh dari parameter.

Dukungan untuk format lain seperti PNG atau WebP belum tersedia.