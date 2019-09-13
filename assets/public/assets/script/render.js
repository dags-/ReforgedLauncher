const Render = {
    el: function(type, props, content) {
        var element = document.createElement(type);
        Render.props(element, props);
        Render.children(element, content);
        return element;
    },
    text: function(string, props) {
        var el = document.createTextNode(string);
        Render.props(el, props);
        return el;
    },
    props: function(el, props) {
        if (!props) {
            return
        }
        Utils.eachObj(props, function (k, v) {
            if (k === "events") {
                return;
            }
            if (k in el) {
                el[k] = v;
            } else {
                el.setAttribute(k, v);
            }
        });
        Utils.eachObj(props["events"], function(k, v) {
            el.addEventListener(k, v.bind(el));
        });
    },
    children: function(el, content) {
        Utils.eachArr(content, function(v) {
            el.appendChild(v);
        });
    },
    redraw: function(el) {
        el.style.visibility = "hidden";
        el.style.visibility = "visible";
    }
};

const Utils = {
    clear: function(el) {
        if (typeof el === "string") {
            el = document.getElementById(el);
        }
        while (el.lastChild) {
            el.removeChild(el.lastChild);
        }
        return el;
    },
    eachArr: function (arr, fn) {
        if (arr) {
            for (var i = 0; i < arr.length; i++) {
                fn(arr[i], i);
            }
        }
    },
    eachObj: function (obj, fn) {
        if (obj) {
            for (var k in obj) {
                if (obj.hasOwnProperty(k)) {
                    fn(k, obj[k]);
                }
            }
        }
    }
};