# 📡 PB Freaking Fast Fetcher 🧑‍💻

This script fetches headers and body content from a list of URLs. It supports concurrent processing of multiple URLs and saves the retrieved data (headers and/or body) into a structured folder format.
This script is an update to the work of [TomNomNom and fff](https://github.com/tomnomnom/fff)

## 🚀 Features
- **Public IP Display**: Fetch and display the public IP of the machine running the script (to avoid nasty surprises).
- **Fetch Headers**: Retrieve HTTP headers for a list of URLs.
- **Fetch Body**: Retrieve the body content of URLs.
- **Parallel Processing**: Process multiple URLs concurrently for faster execution.
- **Output Directory**: Saves the retrieved data into a specified directory structure.
- **SSL Verification Skipped**: Makes requests even if SSL certificates are invalid or missing.

## 🛠️ Requirements
- Go 1.18 or higher

## 💻 Installation

1. Clone the repository or download the script file.
2. Compile the Go program:
   ```bash
   chmod +x build.sh
   ./build.sh 

## 🎯 Usage
Command-Line Arguments:
-f <file>: Path to a file containing a list of URLs (one per line). If not provided, URLs will be read from standard input.
-o <dir>: Output directory for saved files (default: output).
-h: Fetch only HTTP headers for each URL.
-b: Fetch only the body content for each URL.
-n <num>: Number of parallel jobs to process URLs (default: 4).

Example Command:
```
./pbfff -f urls.txt -o ./output -n 20
```
This will:

- Read URLs from urls.txt.
- Fetch body and headers.
- Save the results in the output directory.
- Use 20 parallel jobs to speed up processing.

### Example Input File (urls.txt):
```
http://example.com
https://api.example.com
http://example.org
```
### Example Output:

```
Public IP: 203.0.113.45
[203.0.113.45] Successfully fetched header for http://example.com
[203.0.113.45] Successfully fetched body for https://api.example.com
Finished processing all URLs.
```
## Output Directory Structure:
```
output/
├── example.com/
│   ├── header.txt
│   ├── body.txt
├── api.example.com/
│   ├── header.txt
│   ├── body.txt
```
## ⚙️ How It Works
Public IP Fetching: The script fetches the public IP of the machine running the script from an external service (https://api.ipify.org).
Header and Body Fetching: Depending on the flags passed (-h or -b), it fetches either the HTTP headers or the body (or both) for each URL.
Parallel Processing: The script uses Goroutines to process multiple URLs concurrently, speeding up the fetching process.
Directory Structure: For each URL, a directory is created under the specified output directory. The HTTP headers and/or body content are saved as header.txt and body.txt inside the respective directories.


### 💡 Helpful Notes
The script supports SSL/TLS verification skipping by default, meaning it can fetch content from HTTPS sites even if they have invalid SSL certificates.
The output directory structure is automatically created for each URL, so the data is organized.

### 📝 License
This script is open-source and released under the MIT License.

🙏 Acknowledgements
@tomnonom
@CyberSamurai
