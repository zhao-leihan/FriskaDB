# RayhanDB Query Demo 🚀

Ini adalah contoh lengkap semua query yang bisa dijalankan di **Rayhan Compass** GUI!

---

## 1️⃣ CREATE TABLE - Bikin Tabel

```sql
RAYCREATE RAYTABLE products (
    id NUMBER,
    name TEXT,
    price NUMBER,
    stock NUMBER,
    category TEXT
);
```

**Hasil**: Tabel `products` dibuat dengan 5 kolom

---

## 2️⃣ INSERT - Masukin Data

```sql
RAYERT RAYINTO products (id, name, price, stock, category) 
RAYVALUES (1, 'Laptop Gaming ROG', 25000000, 3, 'Laptop');
```

```sql
RAYERT RAYINTO products (id, name, price, stock, category) 
RAYVALUES (2, 'Mouse Logitech G Pro', 450000, 15, 'Accessories');
```

```sql
RAYERT RAYINTO products (id, name, price, stock, category) 
RAYVALUES (3, 'Monitor LG UltraGear', 5500000, 5, 'Monitor');
```

**Hasil**: 3 produk berhasil dimasukkan

---

## 3️⃣ SELECT - Ambil Semua Data

```sql
RAYLECT * RAYFROM products;
```

**Hasil**: Semua kolom dan semua baris ditampilkan

---

## 4️⃣ SELECT dengan WHERE - Filter Data

### Filter by Price (di bawah 1 juta)
```sql
RAYLECT name, price RAYFROM products RAYWHERE price BELOW 1000000;
```

### Filter by Category
```sql
RAYLECT * RAYFROM products RAYWHERE category = 'Laptop';
```

### Filter by Stock (lebih dari 10)
```sql
RAYLECT name, stock RAYFROM products RAYWHERE stock ABOVE 10;
```

**Hasil**: Hanya produk yang sesuai kondisi ditampilkan

---

## 5️⃣ UPDATE - Ubah Data

### Update Harga
```sql
RAYDATE products RAYSET price = 400000 RAYWHERE name = 'Mouse Logitech G Pro';
```

### Update Stock
```sql
RAYDATE products RAYSET stock = 20 RAYWHERE id = 2;
```

**Hasil**: Data produk diubah

---

## 6️⃣ DELETE - Hapus Data

### Hapus produk tertentu
```sql
RAYLETE RAYFROM products RAYWHERE id = 1;
```

### Hapus produk mahal (di atas 20 juta)
```sql
RAYLETE RAYFROM products RAYWHERE price ABOVE 20000000;
```

**Hasil**: Produk yang sesuai kondisi dihapus

---

## 7️⃣ DESCRIBE - Lihat Schema Tabel

```sql
RAYC products;
```

**Hasil**: Informasi kolom tabel (nama, tipe data) ditampilkan

---

## 8️⃣ SHOW TABLES - Lihat Semua Tabel

```sql
RAYSHOW RAYTABLES;
```

**Hasil**: List semua tabel di database

---

## 🎯 Operator Perbandingan

| Operator SQL Biasa | RayhanDB Keyword | Contoh |
|-------------------|------------------|--------|
| `>` | `ABOVE` | `price ABOVE 1000000` |
| `<` | `BELOW` | `stock BELOW 5` |
| `>=` | `ATLEAST` | `price ATLEAST 500000` |
| `<=` | `ATMOST` | `stock ATMOST 10` |
| `=` | `=` | `category = 'Laptop'` |
| `!=` | `!=` | `stock != 0` |

---

## 📝 Tips Query di Rayhan Compass GUI

1. **Copy-paste query** dari file ini ke Query Editor
2. **Klik "Run Query"** untuk execute
3. **Hasil** akan muncul di bagian Results
4. **F12** untuk open DevTools dan lihat console logs
5. **Sidebar** akan update otomatis setelah CREATE/DROP table

---

## 🚀 Run Demo Script (Terminal)

Untuk test semua query sekaligus:

```bash
cd c:\Users\Rayhan\Music\RayhanDB
go run examples/demo_queries/main.go
```

Script ini akan:
- ✅ Create tabel products
- ✅ Insert 6 produk
- ✅ SELECT semua & dengan filter
- ✅ UPDATE stock
- ✅ DELETE produk mahal
- ✅ Show semua tables

---

**Selamat mencoba! 🎉**
