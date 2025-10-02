function toggleDir(id) {
  const elem = document.getElementById(id);
  if (!elem) return;
  elem.classList.toggle("hidden");
}