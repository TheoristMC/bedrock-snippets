function refreshTheme() {
    document.documentElement.classList.remove("dark")

    const systemThemeIsDark = window.matchMedia("(prefers-color-scheme: dark)").matches
    const theme = localStorage.getItem("theme") ?? "auto"

    if (theme == "dark" || theme == "auto" && systemThemeIsDark) {
        document.documentElement.classList.add("dark")
    }
}

function updateThemeSelectorIcon() {
    const theme = localStorage.getItem("theme")
    const themeSelectorIcon = document.querySelector("#theme-selector-icon")
    if (theme == "light") {
        themeSelectorIcon.src = "/OcSun2.svg"
    } else if (theme == "dark") {
        themeSelectorIcon.src = "/OcMoon2.svg"
    } else {
        themeSelectorIcon.src = "/OcDevicedesktop2.svg"
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const themeSelector = document.querySelector("#theme-selector")

    themeSelector.value = localStorage.getItem("theme") ?? "auto"
    updateThemeSelectorIcon()

    themeSelector.addEventListener("change", (e) => {
        localStorage.setItem("theme", e.target.value)
        refreshTheme()
        updateThemeSelectorIcon()
    })
})
refreshTheme()