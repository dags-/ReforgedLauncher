function renderModpacks(modpacks) {
    var list = document.getElementById("modpacks");
    while (list.lastChild) {
        list.removeChild(list.lastChild);
    }

    Utils.eachArr(modpacks, function (pack) {
        list.appendChild(modpack(pack));
    });

    document.getElementById("back").addEventListener("click", function() {
        nav("/home");
    });

    document.getElementById("add").addEventListener("click", function() {
        document.body.appendChild(inputDialog("", function(text) {
            post("/api/modpacks", JSON.stringify(text), function(result) {
                if (result.success) {
                    nav("/modpacks")
                } else {
                    alert(result["data"])
                }
            });
        }));
    });
}

function modpack(pack) {
    return Render.el("div", {class: "modpack-container theme-row"}, [
        Render.el("img", {src: pack["icon"], class: "modpack-img"}),
        Render.el("div", {class: "modpack-info"}, [Render.text(pack["title"] || pack["name"])]),
        Render.el("div", {class: "modpack-buttons"}, [
            Render.el("input", {
                type: "button", class: "primary", value: "Install", events: {
                    click: function () {
                        document.body.appendChild(inputDialog(pack["title"] || pack["name"], function(text) {
                            post("/api/install/" + pack["repo"], JSON.stringify({name: text}));
                        }));
                    }
                }
            })
        ])
    ]);
}

function inputDialog(value, callback) {
    return Render.el("div", {id: "dialog", class: "input-dialog"}, [
        Render.el("div", {class: "input-container"}, [
            Render.el("input", {type: "text", id: "dialog-input", value: value}),
            Render.el("div", {}, [
                Render.el("input", {
                    type: "button", class: "secondary", value: "Cancel", events: {
                        click: function () {
                            var el = document.getElementById("dialog");
                            el.parentElement.removeChild(el);
                        }
                    }
                }),
                Render.el("input", {
                    type: "button", class: "primary", value: "Confirm", events: {
                        click: function () {
                            var text = document.getElementById("dialog-input").value;
                            callback(text);
                        }
                    }
                }),
            ])
        ])
    ]);
}