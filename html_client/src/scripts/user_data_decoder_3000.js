window.user_data_decoder_3000 = function(){
    // Utility function to try decompressing with pako (gzip)
    function tryDecompress(data) {
        try {
            return pako.inflate(data, {to: 'string'});
        } catch (e) {
            // Return original data if decompression fails
            return data;
        }
    }

    function stringToUint8Array(str) {
        const len = str.length;
        const bytes = new Uint8Array(len);
        for (let i = 0; i < len; i++) {
            bytes[i] = str.charCodeAt(i);
        }
        return bytes;
    }

    function tryParseMultipartMime(mimeContent) {
        try {
            // Assuming mimeContent is a string containing the MIME message.

            // Find the boundary string in the MIME content.
            const boundaryPrefix = 'boundary=';
            const startIndex = mimeContent.indexOf(boundaryPrefix);
            if (startIndex === -1) {
                throw new Error("Boundary not found");
            }

            let endIndex = mimeContent.indexOf('\n', startIndex);
            if (endIndex === -1) endIndex = mimeContent.length;
            let boundary = mimeContent.substring(startIndex + boundaryPrefix.length, endIndex).trim();

            // Handling case where boundary is surrounded by quotes
            if (boundary.startsWith('"') && boundary.endsWith('"')) {
                boundary = boundary.substring(1, boundary.length - 1);
            }

            // Split the content based on the boundary.
            const parts = mimeContent.split(`--${boundary}`);
            const parsedParts = [];
            for (let i = 1; i < parts.length - 1; i++) { // Skip the first and last part
                const part = parts[i].trim();
                if (part === "--") continue; // Skip the closing boundary marker

                // Further processing of each part can be done here.
                parsedParts.push(part); // For simplicity, just adding the raw part
            }

            return parsedParts.length ? parsedParts : mimeContent;
        } catch (e) {
            // console.log("Error parsing MIME parts:", e.message);
            // Return original data if parsing fails
            return mimeContent;
        }
    }

// Placeholder for checking if a string is base64 encoded
    function isBase64Encoded(str) {
        try {
            return btoa(atob(str)) === str;
        } catch (err) {
            return false;
        }
    }

    function parseNonMimeContent(content){
        if (isBase64Encoded(content)) {
            // Decode base64
            content = atob(content);
        }
        content = tryDecompress(content);
        try {
            content = new TextDecoder("utf-8").decode(content);
        }catch(err){}
        if (content.includes("#cloud-config")) {
            // Process as cloud-init YAML content
            return processCloudInit(content);
        } else  {
            // Handle shell scripts
            const filename = "userdata";
            return [{ path: filename, content: content }];
        } // Add other content types if needed
    }

    function parseMimePart(partContent) {
        // Split headers and body
        // Normalize line endings to LF (\n) and then split
        const normalizedContent = partContent.replace(/\r\n/g, '\n');

        // Splitting into sections based on two consecutive newlines,
        // which works regardless of whether it was originally LF or CRLF
        const sections = normalizedContent.split(/\n\n+/);

        if (sections.length < 2) {
            console.error('Failed to parse MIME part: Headers and body are not clearly separated.');
            return { headersPart: '', bodyPart: '' };
        }

        // The first section is headers, and the rest is considered as the body.
        // This handles cases where the body might start with what looks like a separator.
        const headersPart = sections.shift();
        let bodyPart = sections.join('\n\n'); // Rejoin the rest in case the body itself contained \n\n

        const headers = headersPart.split('\n').reduce((acc, current) => {
            const [key, value] = current.split(':', 2).map(s => s.trim());
            acc[key.toLowerCase()] = value; // Use lowercase for header keys for easier matching
            return acc;
        }, {});

        // Check for base64 encoding
        if (headers['content-transfer-encoding'] === 'base64') {
            // Decode base64 content
            bodyPart = atob(bodyPart.trim());
        }
        // Convert to Uint8Array for potential gzip decompression
        let contentArray = stringToUint8Array(bodyPart);

        // Try decompressing (in case it's gzipped), or use as is
        contentArray = tryDecompress(contentArray);

        // Convert Uint8Array back to string
        const content = new TextDecoder("utf-8").decode(contentArray);

        // Determine the type of content based on Content-Type header
        if (headers['content-type']) {
            if (headers['content-type'].includes('cloud-config') ||
                content.includes("#cloud-config")) {
                // Process as cloud-init YAML content
                return processCloudInit(content);
            } else  {
                // Handle shell scripts
                const filename = headers['content-disposition'] ? headers['content-disposition'].split('filename="')[1].split('"')[0] : crypto.randomUUID();
                return [{ path: filename, content }];
            } // Add other content types if needed
        }


        // Return raw body if not base64 encoded
        return [{ path: "userdata", content: bodyPart }];
    }


// Function to process cloud-init YAML content
    function processCloudInit(yamlContent) {
        let files = [];
        try {
            // Assuming a simple YAML parse function is available
            const cloudInitConfig = jsyaml.load(yamlContent);

            // Check for write_files directive or similar in cloud-init config
            if (cloudInitConfig.write_files) {

                cloudInitConfig.write_files.forEach(file => {
                    let content = file.content;
                    if (isBase64Encoded(content)) {
                        // Decode base64
                        content = atob(content);
                        // Convert to Uint8Array for potential gzip decompression
                        content = stringToUint8Array(content);
                        // Try decompressing (in case it's gzipped), or use as is
                        content = tryDecompress(content);

                    }
                    files.push({ path: file.path, content: content });
                });
            }
        } catch (e) {
            console.error("Failed to process cloud-init content:", e);
        }
        return files.length ? files : [{ path: "userdata", content: yamlContent }];
    }


    return {
        tryDecompress,
        stringToUint8Array,
        tryParseMultipartMime,
        parseMimePart,
        parseNonMimeContent
    }
}
