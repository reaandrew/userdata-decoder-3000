
window.user_data_decoder_3000_ui = function(){
    function attach(config){
        let decoder3000 = user_data_decoder_3000()
        let button = document.getElementById(config.button_id)
        button.addEventListener("click", () => {
            input_data = document.getElementById(config.input_id).value
            let value = atob(input_data)
            value = decoder3000.stringToUint8Array(value);
            value = decoder3000.tryDecompress(value)
            let rawDecoded = value;
            try {
                rawDecoded = new TextDecoder("utf-8").decode(rawDecoded);
            }catch(err){}
            value = decoder3000.tryParseMultipartMime(value);

            let fileList = document.getElementById(config.file_list_id);
            let file_content = document.getElementById(config.file_content_id)
            fileList.innerHTML=""
            file_content.innerHTML=""
            let processedFiles = [];
            processedFiles.push({path: "raw_base64decoded", content:rawDecoded});
            if (Array.isArray(value) && value.length > 0) {
                value.forEach(partContent => {
                    let file_parts = decoder3000.parseMimePart(partContent);
                    processedFiles = processedFiles.concat(file_parts);
                });
            }else {
                let files = decoder3000.parseNonMimeContent(value);
                processedFiles = processedFiles.concat(files);
            }

            processedFiles.forEach(function(obj, index) {
                let li = document.createElement("li");
                li.className = "list-group-item";
                li.role = "button"
                if (index == 0){
                    li.classList.add("active")
                    li.classList.add("link")
                    file_content.innerHTML = obj.content;
                }
                // Optionally, add the 'blue' class to all items or based on a condition
                li.textContent = obj.path;

                // Store the object's contents in a custom attribute, e.g., 'data-info'
                // Note: Custom attributes should be all lowercase
                li.setAttribute("data-info", obj.content);

                // Add click event listener to list item
                li.addEventListener('click', function() {
                    // Remove 'active' class from all list items
                    document.querySelectorAll('.list-group-item').forEach(function(el) {
                        el.classList.remove('active');
                    });

                    // Add 'active' class to clicked list item
                    this.classList.add('active');

                    // Set the file_content to the data-info attribute of the clicked list item

                    file_content.innerHTML = this.getAttribute('data-info');
                });

                fileList.appendChild(li);
            });

        })
    }

    return {
        attach
    }
}