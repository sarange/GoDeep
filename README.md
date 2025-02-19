# GoDeep

GoDeep is an steganographic tool designed to embed and extract files within WAV audio using least significant bit (LSB) encoding. Unlike existing proprietary solutions like DeepSound, which are closed-source and limited to Windows, GoDeep is open-source and cross-platform, allowing greater accessibility and flexibility for users. It provides a secure and efficient way to conceal data within audio files while maintaining the integrity and quality of the original sound. The tool also offers optional AES-GMC encryption.

## Embedded File Structure

The embedded file follows a structured format to ensure accurate extraction and decryption when needed. This format maintains data integrity and facilitates seamless retrieval of hidden content.

| Type               | Size in Bytes                           | Description                                    |
| ------------------ | --------------------------------------- | ---------------------------------------------- |
| Magic Bytes        | 3 (`byte[3]`) ("GDP" -> `\x47\x44\x50`) | Identifies the embedded file format.           |
| Encryption         | 1 (`bool`)                              | Indicates whether encryption is enabled.       |
| Size of Nonce      | 1 (`uint8`)                             | Specifies the length of the nonce.             |
| Nonce              | 0-255 (based on `Size of Nonce`)        | Random value for encryption (if enabled).      |
| Size of Ciphertext | 8 (`uint64`)                            | Length of the encrypted data or plaintext.     |
| Ciphertext         | 0 - 18,446,744,073,709,551,615 bytes    | The actual embedded data (encrypted or plain). |
| ----               | ----                                    | ----                                           |
| **Total**          | `13 + len(nonce) + len(ciphertext)`     | The total size of the embedded file.           |

## Features
- **LSB Encoding**: Efficiently conceals data within the least significant bits of PCM audio without introducing audible distortion.
- **Lossless Extraction**: Ensures accurate retrieval of hidden files, preserving data integrity even after multiple extractions.
- **Encryption Support**: Offers optional AES encryption for added security, protecting sensitive data from unauthorized access.
- **Flexible File Size Handling**: Supports embedding both small and large files while maintaining the playability of the audio.
- **Cross-Platform Compatibility**: Works on various operating systems, making it accessible to a broad range of users.
- **Optimized Performance**: Utilizes efficient encoding and decoding algorithms to minimize processing time and maximize accuracy.
- **User-Friendly CLI**: Provides an intuitive command-line interface for straightforward embedding and extraction of files.

## Usage
### Embedding a File
Easily conceal a file within a WAV audio file:
```sh
./godeep embed -i input.wav -f secret.txt -o output.wav -p "GoDeep"
```

### Extracting an Embedded File
Retrieve hidden data from a WAV file:
```sh
./godeep extract -c output.wav -o extracted.txt -p "GoDeep"
```

### Embedding a File without encryption
Enhance security by encrypting the file before embedding:
```sh
./godeep embed -i input.wav -f secret.txt -o output.wav --noencryption
```

## License
GoDeep is released under the **GNU General Public License v3.0**, allowing free usage, modification, and distribution with proper attribution. We encourage community contributions to improve and expand the tool.

## Contributions
We welcome contributions! Whether it's bug fixes, feature enhancements, or optimizations, your input helps improve GoDeep. Feel free to submit pull requests or open issues on our GitHub repository.

## ToDo
- Add the ability to chose the depth of the embeding
- Add a check command
