window.onload = function () {
    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get('id');

    var conn;
    var contentReady = false;
    var readyState = false;
    var editorChange = false;

    var clientContainer = document.getElementById("clients");
    var btnSave = document.getElementById("save");
    var btnReady = document.getElementById("ready");
    var btnUnready = document.getElementById("unready");
    var btnDisconnect = document.getElementById("disconnect");

    var editor = ace.edit("editor");
    editor.setTheme("ace/theme/GitHub");
    editor.session.setUseSoftTabs(true);
    editor.resize();
    editor.setOptions({
        fontSize: "12pt"
    });


    btnSave.addEventListener("click", () => {
        var save = {
            "type": "conn-save",
        };

        conn.send(JSON.stringify(save));
    });

    btnDisconnect.addEventListener("click", () => {
        var disconnect = {
            "type": "conn-disconnect",
        };
        
        conn.send(JSON.stringify(disconnect));
    });

    btnReady.addEventListener("click", () => {
        btnReady.style.display = "none";
        btnReady.disabled = "disabled";

        btnUnready.style.display = "block";
        btnUnready.disabled = "";
        
        editor.setReadOnly(true);

        var ready = {
            "type": "conn-ready",
        };

        readyState = true;
        
        conn.send(JSON.stringify(ready));
    });

    btnUnready.addEventListener("click", () => {
        btnReady.style.display = "block";
        btnReady.disabled = "";

        btnUnready.style.display = "none";
        btnUnready.disabled = "disabled";

        editor.setReadOnly(false);

        var ready = {
            "type": "conn-unready",
        };

        readyState = false;
        
        conn.send(JSON.stringify(ready));
    });


    editor.on('change', function(delta) {
        if (contentReady == false){
            return;
        }

        if (editorChange == true){
            return;
        }

        var change = {
            "type": "conn-text-change",
            "data": delta
        };

        conn.send(JSON.stringify(change));
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws?id=" + id);

        conn.onclose = function () {
            editor.setReadOnly(true);
            clientContainer.innerHTML = "";
            
            location.reload();
        };

        conn.onmessage = function (evt) {
            var update = JSON.parse(evt.data);

            // console.log("event received: ", evt);

            switch (update.type) {
                case "conn-connected":
                    editor.setReadOnly(true);
                    editor.setValue(update.data);

                    editor.session.setMode("ace/mode/" + update.fileMeta.extension);

                    editor.setReadOnly(false);
                    contentReady = true;
                    break;
                case "conn-disconnected":
                    console.log("connection closed");
                    break;
                case "clients-text-change":
                    var updateContents = JSON.parse(update.data);

                    editorChange = true;
                    editor.session.doc.applyDelta(updateContents.data);
                    editorChange = false;
                    break;
                case "clients-connected":
                    showAlert(update.client + " connected", "success");
                    break;
                case "clients-disconnected":
                    showAlert(update.client + " disconnected", "warning");
                    break;
                case "server-file-saved":
                    clearAlerts();

                    showAlert("File saved!", "success");
                    readyState = false;

                    hideUnreadyButton(btnUnready, editor);
                    break;
                case "server-file-not-ready":
                    showAlert("File not ready!", "danger");
                    break;
                case "server-file-not-saved":
                    showAlert("File not saved!", "danger");
                    break;
                case "conn-not-ready":
                    showAlert("Connection not ready!", "danger");
                    break;
                case "conn-not-unready":
                    showAlert("Connection not unready!", "danger");
                    break;
            }

            refreshClients(clientContainer, btnSave, btnReady, btnUnready, readyState, update.clients);
        };
    } else {
        console.log("Your browser does not support WebSockets", "red");
    }
};

function refreshClients(clientContainer, btnSave, btnReady, btnUnready, readyState, clients) {
    clientContainer.innerHTML = "";

    if (clients.length == 1) {
        btnSave.disabled = "";
        btnSave.style.display = "block";

        hideReadyButtons(btnReady, btnUnready);
    } else {
        btnSave.disabled = "disabled";
        btnSave.style.display = "none";

        if (!readyState){
            showReadyButton(btnReady);
        }
    }

    clients.forEach(function (client) {
        var cl = document.createElement("button");

        cl.innerText = client.username;
        cl.style.color = client.color;
        cl.style.border = "1px solid " + client.color;
        cl.className = 'btn btn-sm btn-light m-1';

        if (client.ready) {
            cl.innerHTML= cl.innerHTML + '<i class="bi bi-check-circle-fill"></i>';
        }
        
        clientContainer.appendChild(cl);
    });
}

function hideReadyButtons(btnReady, btnUnready){
    btnReady.style.display = "none";
    btnReady.disabled = "disabled";

    btnUnready.style.display = "none";
    btnUnready.disabled = "disabled";
}

function showReadyButton(btnReady){
    btnReady.style.display = "block";
    btnReady.disabled = "";
}

function hideUnreadyButton(btnUnready, editor){
    btnUnready.style.display = "none";
    btnUnready.disabled = "disabled";

    editor.setReadOnly(false);
}

function clearAlerts() {
    var alertContainer = document.getElementById("alerts");
    alertContainer.innerHTML = "";
}

function showAlert(text, color) {
    var alertContainer = document.getElementById("alerts");
    alertContainer.innerHTML += `<div class="toast align-items-center text-bg-`+color+` border-0 show" role="alert" aria-live="assertive" aria-atomic="true">\
        <div class="d-flex">\
            <div class="toast-body">\
            `+text+`\
            </div>\
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>\
        </div>\
    </div>`;
}