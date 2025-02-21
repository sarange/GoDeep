GoDeep

GoDeep is a steganographic tool designed to embed and extract files within WAV audio using least significant bit (LSB) encoding. Unlike proprietary solutions like DeepSound, which are closed-source and limited to Windows, GoDeep is open-source and cross-platform, allowing greater accessibility and flexibility for users. It provides a secure and efficient way to conceal data within audio files while maintaining the integrity and quality of the original sound. The tool also offers optional AES-GCM encryption.

## GUI Preview (Planned)

![GUI Screenshot](/images/godeep_gui.png)

## Embedded File Structure

The embedded file follows a structured format to ensure accurate extraction and decryption when needed. This format maintains data integrity and facilitates seamless retrieval of hidden content.

| Type               | Size in Bytes                           | Description                                    |
| ------------------ | --------------------------------------- | ---------------------------------------------- |
| Magic Bytes        | 3 (`byte[3]`) (`GDP` -> `\x47\x44\x50`) | Identifies the embedded file format.           |
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
- **Encryption Support**: Offers optional AES-GCM encryption for added security, protecting sensitive data from unauthorized access.
- **Cross-Platform Compatibility**: Works on Linux, macOS, and Windows.
- **User-Friendly CLI**: Provides an intuitive command-line interface for straightforward embedding and extraction of files.
- **GUI Support**: A graphical user interface to make usage more accessible.
- **Custom Embedding Depth (Planned)**: Users will be able to control the depth of LSB embedding.
- **Integrity Check Command (Planned)**: A command to verify the integrity of embedded files.

## Usage

GoDeep provides both **command-line** and **GUI** options for embedding and extracting files from WAV audio.

---

### **Command-Line Usage**

#### **Embedding a File**
Embed a file into a WAV container:
```sh
godeep embed -i input.wav -c container.wav -o output.wav -p "your_password"
```

- `-i, --input` → Path to the file you want to embed.
- `-c, --container` → WAV file that will store the hidden data.
- `-o, --output` → Output WAV file containing the embedded data.
- `-p, --password` → Encryption password (unless `--noencryption` is used).

#### **Extracting a File**
Extract hidden data from a WAV file:
```sh
godeep extract -c container.wav -o extracted_file -p "your_password"
```

- `-c, --container` → WAV file that contains the hidden data.
- `-o, --output` → Output file where extracted data will be saved.
- `-p, --password` → Encryption password (if encryption was used).

#### **Embedding a File Without Encryption**
If you want to disable encryption:
```sh
godeep embed -i input.wav -c container.wav -o output.wav --noencryption
```

---

### **Using the GUI**

GoDeep includes a **graphical user interface (GUI)** for users who prefer a visual approach instead of the command line.

#### **Launching the GUI**
To open the GUI, run:
```sh
godeep gui
```

#### **How to Embed a File**
1. Run `godeep gui` to open the graphical interface.
2. Click **Embed** to start hiding a file inside a WAV.
3. Select the **container WAV file**, the **file to embed**, and specify an **output file**.
4. Enter a password (if encryption is enabled).
5. Click **Run** to embed the file.

#### **How to Extract a File**
1. Open the GUI by running `godeep gui`.
2. Click **Extract** to retrieve hidden data from a WAV file.
3. Select the **container WAV file** and specify an **output file**.
4. Enter a password (if encryption was used).
5. Click **Run** to extract the file.

## Testing

GoDeep includes automated tests to ensure reliability and correctness.

To run the tests, use:
```sh
go test -v
```

This will execute all tests in the project, including unit and integration tests.

Contributions that include new features or bug fixes should also include relevant test cases.

## License

GoDeep is released under the **GNU General Public License v3.0 (GPL-3.0)**, allowing free use, modification, and distribution, as long as proper attribution and license terms are followed. Community contributions are welcome to help improve and expand the tool.

## Contributions

Whether it's bug fixes, feature enhancements, or optimizations, your input helps improve GoDeep. Feel free to submit pull requests or open an issue.

## ToDo

- [x] Implement AES-GCM encryption *(Completed)*
- [x] Add a GUI *(Completed)*
- [ ] Allow users to choose the depth of embedding *(Planned)*
- [ ] Implement an integrity check command *(Planned)*
- [ ] Improve performance optimizations *(Ongoing)*

## Disclaimer

This is a weekend project that I decided to make open-source. While I’ll try to update it when I find the time, please don’t expect regular maintenance or major updates. Contributions are always welcome!
