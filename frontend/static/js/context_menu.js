$(function() {
    $.contextMenu({
        className: 'custom-contextmenu-style',
        selector: '#control-box',
        trigger: 'left',
        callback: function(key, options) {
            switch (key) {
                case "Register":
                    $('#register-modal').modal('show');
                    break
                case "Login":
                    $('#login-modal').modal('show');
                    break
                default:
                    break
            }
        },
        items: {
            "Register": { name: "Register", icon: "fa-pen" },
            "Login": { name: "Login", icon: "fa-lock" },
            "sep1": "---------",
            "quit": { name: "Quit", icon: function($element, key, item) { return 'context-menu-icon context-menu-icon-quit'; } }
        }
    });
});

$.contextMenu({
    className: 'custom-contextmenu-style',
    selector: '.vis-network canvas',
    trigger: 'right',
    callback: function(key, options) {
        switch (key) {
            case "change_color":
                break
            case "change_shape":
                break
            case "change_label":
                break
            default:
                break
        }
    },
    items: {
        "change_shape": {
            name: "Change Shape",
            type: "select",
            options: {
                "": '',
                "circle": 'Circle',
                "ellipse": 'Ellipse',
                "database": 'Database',
                "box": 'Square',
                "triangle": 'Triangle',
                "triangleDown": 'Triangle Down',
                "diamond": 'Diamond',
                "dot": 'Dot',
                "star": 'Star',
                "hexagon": 'Hexagon',
                "text": 'Text',
            }
        },
        "change_label": {
            name: "Change Label",
            type: "text",
        },
        "change_color": { name: "Change Color", icon: "fa-eye-dropper" },
        "sep1": "---------",
        "quit": { name: "Cancel", icon: function($element, key, item) { return 'context-menu-icon context-menu-icon-quit'; } }
    }
});

$(document).ready(function() {
    $('.fa-eye-dropper').colorpicker({
        color: "#ffffff", // Initial color
        defaultPalette: "theme", // Color palette theme
        transparentColor: true // Can select transparent colors
    });
    // Changing node colors
    $('.fa-eye-dropper').on("change.color", function(event, color) {
        selectedNodes = network.getSelectedNodes()
        for (index in selectedNodes) {
            nodeID = selectedNodes[index]
            node = allNodes.get(nodeID)
            if (node != undefined) {
                // Need to handle color in rgba format and default for color picker is hex
                previousColor = parseRGBA(node.color)
                pastAlpha = previousColor[3]
                node.color = hexToRGBA(color, pastAlpha)
                allNodes.update(node)
                allNodes.update(node)
                if (remoteNodes.get(nodeID) != undefined) {
                    remoteNodes.update(node)
                }
                if (agentNodes.get(nodeID) != undefined) {
                    agentNodes.update(node)
                }
            }
        }
    });
    // Changing node shapes
    $('select').change(function() {
        selectedNodes = network.getSelectedNodes()
        newShape = $(this).val()
        for (index in selectedNodes) {
            nodeID = selectedNodes[index]
            node = allNodes.get(nodeID)
            if (node != undefined) {
                node.shape = newShape
                allNodes.update(node)
                if (remoteNodes.get(nodeID) != undefined) {
                    remoteNodes.update(node)
                }
                if (agentNodes.get(nodeID) != undefined) {
                    agentNodes.update(node)
                }
            }
        }
    });
    // Adding node labels/text
    $('input').on('keydown', function(event) {
        if (event.which == 13 && !event.shiftKey) {
            event.preventDefault();
            newLabel = $(this).val()
            selectedNodes = network.getSelectedNodes()
            for (index in selectedNodes) {
                nodeID = selectedNodes[index]
                node = allNodes.get(nodeID)
                if (node != undefined) {
                    node.label = newLabel
                    allNodes.update(node)
                    if (remoteNodes.get(nodeID) != undefined) {
                        node.title = remoteNodeHover(node)
                        remoteNodes.update(node)
                    }
                    if (agentNodes.get(nodeID) != undefined) {
                        node.title = agentNodeHover(node)
                        agentNodes.update(node)
                    }
                }
            }
        }
    });
});