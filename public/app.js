(function (d, axios) {
    "use strict";
    var inputFile = d.querySelector("#inputFile");
    var divNotification = d.querySelector("#alert");
    var buttonMix = d.querySelector("#button-mix");
    var buttonInput = d.querySelector("#button-input");
    var mixes = d.querySelector("#mixes");
    var inputs = d.querySelector("#inputs");


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
            response.data.Files.map(function (x){
                inputs.innerHTML += "<li> <div style='background-color:darkseagreen'>"
                inputs.innerHTML +=  "input- <a href=" + "'input/"+ x +"'" + ">" + x + "</a>"
                //inputs.innerHTML += "<button> Mix </button>"
                inputs.innerHTML += " </div></li>"
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

            response.data.Files.map(function (x){
                mixes.innerHTML += "<li> <div>"
                mixes.innerHTML +=  "mix <a href=" + "'output/"+ x +"'" + ">" + x + "</a>"
                mixes.innerHTML += " </div></li>"
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
            .then(onResponse)
            .catch(onResponse);
    }

    function onResponse(response) {
        var className = (response.status !== 400) ? "success" : "error";
        divNotification.innerHTML = response.data;
        divNotification.classList.add(className);
        setTimeout(function() {
            divNotification.classList.remove(className);
        }, 3000);
    }




})(document, axios)