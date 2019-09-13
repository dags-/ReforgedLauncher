function renderSettings(settings) {
    var header = Utils.clear("header");
    header.innerText = settings["name"] + " (" + settings["mod_pack"]["repo"] + ")";

    var options = Utils.clear("options");
    options.appendChild(gameDir(settings));
    options.appendChild(coverImage(settings));

    Utils.eachObj(settings["options"], function (name, option) {
        options.appendChild(optionalMod(option, settings["options"]));
    });

    document.getElementById("back").addEventListener("click", function () {
        var path = "/api/instance/" + settings["name"];
        var body = JSON.stringify(settings);
        post(path, body, function () {
            nav("../home");
        });
    });

    document.getElementById("delete").addEventListener("click", function () {
        var path = "/api/instance/" + settings["name"];
        del(path, function () {
            nav("/home");
        });
    });

    document.getElementById("launch").addEventListener("click", function () {
        var path = "/api/instance/" + settings["name"];
        var body = JSON.stringify(settings);
        post(path, body, function () {
            get("/api/run/" + settings["name"])
        });
    });

    refresh(settings["options"]);

    Render.redraw(document.body);
}

function gameDir(settings) {
    return Render.el("label", {for: "game_dir", class: "option theme-row"}, [
        Render.el("div", {class: "expand"}, [
            Render.el("div", {class: "text"}, [
                Render.text("Game Dir")
            ]),
            Render.el("div", {class: "row"}, [
                Render.el("input", {
                    id: "gameDir", type: "text", value: settings["game_dir"], events: {
                        change: function () {
                            settings["game_dir"] = this.value;
                        }
                    }
                }),
                Render.el("input", {
                    type: "button", class: "secondary open-button", value: "Open", events: {
                        click: function () {
                            var data = {name: settings["game_dir"]};
                            post("/api/open/folder", JSON.stringify(data));
                        }
                    }
                })
            ])
        ])
    ])
}

function coverImage(settings) {
    return Render.el("label", {for: "coverImage", class: "option theme-row"}, [
        Render.el("div", {class: "expand"}, [
            Render.el("div", {class: "text"}, [
                Render.text("Cover Image")
            ]),
            Render.el("div", {class: "row"}, [
                Render.el("input", {
                    id: "coverImage", type: "text", value: settings["image"], events: {
                        change: function () {
                            settings["image"] = this.value;
                        }
                    }
                })
            ])
        ])
    ])
}

function optionalMod(option, options) {
    return Render.el("label", {for: option.name, class: "option theme-row"}, [
        Render.el("div", {class: "expand"}, [
            Render.el("div", {class: "text"}, [Render.text(option.name)]),
            Render.el("div", {class: "sub-text"}, [Render.text(option.description || "-")])
        ]),
        Render.el("div", {class: "switch-container"}, [
            Render.el("input", {
                id: option.name,
                name: option.name,
                type: "checkbox",
                checked: option.enabled,
                events: {
                    change: function () {
                        refresh(options);
                    }
                }
            }),
            Render.el("span", {class: "switch"})
        ]),
    ]);
}

function refresh(options) {
    Utils.eachObj(options, function (name, option) {
        update(option, options);
    });
}

function update(option, options) {
    var enabled = true;
    var el = document.getElementById(option["name"]);
    Utils.eachArr(option["dependencies"], function (name) {
        var dep = options[name];
        if (!update(dep, options)) {
            enabled = false;
        }
    });
    el.enabled = enabled;
    el.checked = enabled && el.checked;
    option["enabled"] = el.checked;
    if (enabled) {
        el.parentElement.parentElement.classList.remove("disabled");
    } else {
        el.parentElement.parentElement.classList.add("disabled");
    }
    return el.checked;
}