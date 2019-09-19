function renderModpacks(packs) {
    var list = document.getElementById("instances");

    while (list.lastChild) {
        list.removeChild(list.lastChild);
    }

    Utils.eachArr(packs, function (pack) {
        list.appendChild(instance(pack));
    });

    list.appendChild(create());
}

function loadScale(config) {
    document.getElementById("scale").value = config["home_scale"];

    var handle;
    document.getElementById("scale").addEventListener("change", function () {
        var scale = this.value;
        clearTimeout(handle);

        handle = setTimeout(function() {
            config["home_scale"] = parseInt(scale);
            post("/api/home/scale", JSON.stringify(config));
        }, 500);

        setScale(scale);
    });

    setScale(config["home_scale"]);
}

function instance(data) {
    return Render.el("div", {class: "instance-container"}, [
        Render.el("div", {
            class: "instance theme-contrast", events: {
                click: function () {
                    nav("/instance#" + data["name"])
                }
            }
        }, [
            Render.el("img", {
                class: "instance-img", src: data["cover1"] || data["cover2"], events: {
                    error: function () {
                        this.src = "../assets/image/banner.jpg";
                    }
                }
            }),
            Render.el("div", {class: "instance-info"}, [
                Render.el("div", {class: "text"}, [Render.text(data["name"])]),
                Render.el("div", {class: "sub-text"}, [Render.text("(" + data["pack"] + ")")]),
            ]),
            Render.el("div", {class: "instance-overlay"})
        ])
    ]);
}

function setScale(scale) {
    scale = 10 - scale;
    var width = ((1 / scale) * 100) + "%";
    var height = (40 / (scale / 3)) + "vw";
    var els = document.getElementsByClassName("instance-container");
    Utils.eachArr(els, function (el) {
        el.style.flexBasis = width;
        el.style.fontSize = width;
        el.style.height = height;
    });
}

function create() {
    return Render.el("div", {class: "instance-container"}, [
        Render.el("div", {
            class: "instance instance-add-container theme-contrast", events: {
                click: function () {
                    nav("/modpacks")
                }
            }
        }, [
            Render.el("div", {class: "instance-add"}, [Render.text("+")]),
            Render.el("div", {class: "instance-overlay"})
        ])
    ]);
}