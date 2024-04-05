document.addEventListener("DOMContentLoaded", function() {

    let ui = user_data_decoder_3000_ui()
    ui.attach({
        button_id: 'decode_button',
        input_id: 'base64_input',
        file_list_id: 'file_list',
        file_content_id: 'file_content'
    })
});