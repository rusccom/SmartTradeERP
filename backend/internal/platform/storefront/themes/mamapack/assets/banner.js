const slides = document.getElementById("banner-slides");
const controls = document.querySelector(".banner-controls");
const dots = Array.from(controls.querySelectorAll("[data-slide]"));

function currentSlide() {
  return Math.round(slides.scrollLeft / slides.clientWidth);
}

function indicateSlide() {
  const current = dots[currentSlide()];
  for (const dot of dots) {
    if (dot === current) dot.setAttribute("aria-current", "true");
    else dot.removeAttribute("aria-current");
  }
}

function selectSlide(event) {
  const button = event.target.closest("button");
  if (!button) return;
  const next = button.hasAttribute("data-slide") ? Number(button.dataset.slide) : currentSlide() + Number(button.dataset.shift);
  const index = (next + dots.length) % dots.length;
  slides.scrollTo({ left: slides.children[index].offsetLeft });
}

controls.addEventListener("click", selectSlide);
slides.addEventListener("scroll", indicateSlide, { passive: true });
controls.hidden = false;
