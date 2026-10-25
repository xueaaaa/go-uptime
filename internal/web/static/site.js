document.querySelectorAll("[data-dialog]").forEach(button => {
    button.addEventListener("click", () => {
        document.getElementById(button.dataset.dialog).showModal();
    });
});

document.querySelector(".site-dialog[data-open]")?.showModal();
