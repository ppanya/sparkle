##

- Provider merge กันได้ trust email จาก Line, Facebook
- Email สามารถส่ง template เข้ามาได้
- รูป profile ใช้ตั้งต้นมาจาก provider
- ถ้าไม่มี primary display picture จะใช้ latest provider profile picture




## MongoDB replica for testing
- ใช้ `run-rs` ในการทำ replica set ใน local ก่อนที่จะ integration test mongodb
```
$ run-rs
```