# 🌍 // GEO_STEGO

> **"A picture is worth a thousand words."** > *With GeoStego, that is literally true.*

![Go](https://img.shields.io/badge/backend-Go_1.23-00ADD8?style=for-the-badge&logo=go)
![Angular](https://img.shields.io/badge/frontend-Angular_19_(Zoneless)-DD0031?style=for-the-badge&logo=angular)
![Docker](https://img.shields.io/badge/deploy-Docker-2496ED?style=for-the-badge&logo=docker)
![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge)

---

## 🕵️‍♂️ The Concept
**GeoStego** is a cryptographic tool that fuses **Steganography** (hiding data in pixels) with **Geofencing** (TODO).

Most encryption protects "What" (the data). GeoStego protects "Where".
1.  **Lock:** You hide a message inside a standard PNG image.
2.  **Geotag:** The encryption key is derived from the GPS coordinates provided.
3.  **Unlock:** The recipient can only decode the message if they are nearby the coordinates provided (TODO : Radius needs to be defined).

### 🧠 Technical Capacity

GeoStego uses **LSB (Least Significant Bit)** encoding across the **RGB** channels. This allows for significant data storage without visible distortion.

**The Formula:**
$$C = \frac{W \times H \times 3}{8} - 4$$

Where:
* $W, H$: Image Dimensions
* $3$: Usage of Red, Green, and Blue channels
* $4$: Bytes reserved for message length header (uint32)

### 📏 The "1000 Words" Benchmark

> **"A picture is worth a thousand words."**

In GeoStego, this is mathematically verified.
To store exactly **1,000 words** (approx. 6KB of text), you don't need a 4K wallpaper.

You only need a **127x127 pixel** icon.

| Message Size | Required Image Size (Min) |
| :--- | :--- |
| **1 Word** (Password) | 4x4 px |
| **1 Tweet** (280 chars) | 28x28 px |
| **1,000 Words** (Essay) | **127x127 px** |
