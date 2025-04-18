const snippetContent = document.querySelector("#snippet-content")
const contentType = snippetContent?.getAttribute("data-content-type")

const copyButton = document.querySelector("#copy-content")

function init() {
    if (copyButton === null) return
    copyButton.onclick = () => {
        if (snippetContent === null) {
            displayToast("Failed to copy to the clipboard.")
            throw Error("Snippet content is missing")
        }

        if (contentType === "image") {
            copyImageContent()
        } else if (contentType == "text") {
            copyTextContent()
        } else {
            displayToast("Failed to copy to the clipboard.")
            throw Error("Unexpected content type: " + contentType)
        }
    }
}

init()

function copyImageContent() {
    try {
        // create a new image with the same content, but sized correctly
        const newImage = document.createElement("img")
        newImage.src = snippetContent.src

        const canvas = document.createElement("canvas")
        canvas.width = newImage.width
        canvas.height = newImage.height

        const ctx = canvas.getContext("2d")
        ctx.drawImage(newImage, 0, 0, newImage.width, newImage.height)

        canvas.toBlob((blob) => {
            navigator.clipboard.write([
                new ClipboardItem({
                    'image/png': blob
                })
            ])
        }, "image/png")
        displayToast("Copied image to clipboard.")
    } catch (error) {
        displayToast("Failed to copy to clipboard.")
        
        throw error
    }
}

function copyTextContent() {
    try {
        const textContent = snippetContent?.getAttribute("data-content-text")

        navigator.clipboard.write([
            new ClipboardItem({
                'text/plain': textContent
            })
        ])
        displayToast("Copied text to clipboard.")
    } catch (error) {
        displayToast("Failed to copy to clipboard.")
        
        throw error
    }
}

const toast = document.querySelector("#toast")
let hideToastTimeout;
function displayToast(string) {
    toast.classList.remove("animate-toast")
    toast.innerText = string
    toast.style.display = undefined
    clearTimeout(hideToastTimeout)

    setTimeout(() => {
        toast.classList.add("animate-toast")
    }, 10)
    hideToastTimeout = setTimeout(() => {
        toast.classList.remove("animate-toast")
    }, 3000)
}