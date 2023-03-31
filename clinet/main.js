var components = {
    message: `<div>
    <p class="message-content">${content}</p>
    <div class="message-about">
        <p class="message-author">${author}</p>
        <p class="subheading">${timestamp}</p>
    </div>
</div>`
}

function addComponent(element) {
    var inner = components[element.tagName]
    var attributeDefinitions = ""
    for (let attr of element.attributes) {
        console.log(attribute)
        attributeDefinitions += ", "
    }
    node.innerHTML = eval("var")
}

var c = document.createElement("message")
c.setAttribute("content", "hi this is a message")
c.setAttribute("author", "JohnS")
c.setAttribute("timestamp", "23/23/23")

addComponent(c)

