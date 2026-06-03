export function scrollContainerToTop(
  container: HTMLElement,
  behavior: ScrollBehavior = "smooth",
) {
  container.scrollTo({ top: 0, behavior });
}
