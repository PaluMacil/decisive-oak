function loadJSON(path, success, error) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                if (success)
                    success(JSON.parse(xhr.responseText));
            } else {
                if (error)
                    error(xhr);
            }
        }
    };
    xhr.open("GET", path, true);
    xhr.send();
}

function convertTree(root) {
    let chart_config = {};
    chart_config.chart = {
        container: "#oak-tree",

        nodeAlign: "BOTTOM",

        connectors: {
            type: 'step'
        },
        node: {
            HTMLclass: 'nodeExample1'
        }
    };
    chart_config.nodeStructure = convertNode(root);

    new Treant(chart_config);
}

function convertNode(node) {
    let newNode = {
        text: {
            filter: node.FilterValue,
            label: node.Label
        },
        HTMLclass: node.Terminal ? 'blue' : 'light-gray'
    };

    if (node.Children && node.Children.length > 0) {
        newNode.children = [];
        node.Children.forEach(child => {
            newNode.children.push(convertNode(child));
        });
    }

    return newNode;
}

loadJSON('/api/list/files',
    function (files) {
        if (files && files.length > 0) {
            showFile(files[0].Filename);
            const btnGroup = document.querySelector('div.btn-group');
            for (const file of files) {
                const buttonEle = document.createElement("button");
                buttonEle.onclick = function () { showFile(file.Filename); };
                buttonEle.innerText = file.Filename;
                btnGroup.appendChild(buttonEle);
            }
        }
    },
    function (xhr) { console.error(xhr); }
);

function showFile(filename) {
    loadJSON(filename,
        function (data) {
            convertTree(data);
        },
        function (xhr) { console.error(xhr); }
    );
}
