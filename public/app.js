(function (d, axios) {
    "use strict";
    var inputFile = d.querySelector("#inputFile");
    var divNotification = d.querySelector("#alert");
    var buttonMix = d.querySelector("#button-mix");
    var buttonInput = d.querySelector("#button-input");
    var mixes = d.querySelector("#mixes");
    var inputs = d.querySelector("#inputs");
    var lastMix = d.querySelector("#lastMix");


    buttonMix.addEventListener("click", buttonMixClick);
    buttonInput.addEventListener("click", buttonInputClick);
    inputFile.addEventListener("change", addFile);

    function buttonInputClick(){
        getJSON("/input")
            .then(onInputResponse)
            .catch(onInputResponse)
    }

    function onInputResponse(response) {
        if (response.status == 200){
            inputs.innerHTML = ""
            response.data.Files.map(function (x){
                inputs.innerHTML += " <div style='background-color:darkseagreen'>"
                inputs.innerHTML +=  "input- <a href=" + "'input/"+ x +"'" + ">" + x + "</a>"
                //inputs.innerHTML += "<button> Mix </button>"
                inputs.innerHTML += "<img height='100' width='200' src='input/" + x+"'/>"
                inputs.innerHTML += " </div>"
            })
        }else{
            alert('error with mixed images')
        }
    }

    function buttonMixClick(){
        console.log("click")
        getJSON("/mixed")
            .then(onMixedResponse)
            .catch(onMixedResponse)
    }

    function onMixedResponse(response) {
        if (response.status == 200){
            mixes.innerHTML = ""
            response.data.Files.map(function (x){
                mixes.innerHTML += "<li> mix <a href=" + "'output/"+ x +"'" + ">" + x + "</a>"
                mixes.innerHTML += " </li>"
            })
        }else{
            alert('error with mixed images')
        }
    }

    function addFile(e) {
        var file = e.target.files[0]
        if(!file){
            return
        }
        upload(file);
    }

    function upload(file) {
        var formData = new FormData()
        formData.append("file", file)
        post("/upload", formData)
            .then(onUploadResponse)
            .catch(setFlashNotification("error uploading"));
    }

    function onUploadResponse(response) {

        var className;
        if (response.status !== 400){
            className = "sucess"
            lastMix.innerHTML = "<img height=500 width='500' src='" + response.data + "'>"
            divNotification.innerHTML = "File uploaded successfully : " + response.data;
        }else{
            className = "error"
            divNotification.innerHTML = response.data;
        }
        
        divNotification.classList.add(className);
        setTimeout(function() {
            divNotification.classList.remove(className);
        }, 3000);
    }

    function setFlashNotification(message){
        var className = "warning"
        divNotification.innerHTML = message;
        divNotification.classList.add(className);
        setTimeout(function() {
            divNotification.classList.remove(className);
        }, 3000);
    }



})(document, axios)