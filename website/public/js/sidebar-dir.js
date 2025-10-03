function toggleDir(id) {
  const elem = document.querySelector(`[data-directory-content="${id}"]`)
  if (!elem) return;
  elem.classList.toggle("hidden");
}