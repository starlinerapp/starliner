export function scrollContainerToTop(
  container: HTMLElement,
  behavior: ScrollBehavior = "smooth",
) {
  container.scrollTo({ top: 0, behavior });
}

export function scrollContainerToSectionStart(
  container: HTMLElement,
  section: HTMLElement,
  behavior: ScrollBehavior = "smooth",
) {
  const containerTop = container.getBoundingClientRect().top;
  const sectionTop = section.getBoundingClientRect().top;
  const nextTop = container.scrollTop + (sectionTop - containerTop);

  container.scrollTo({
    top: Math.max(0, nextTop),
    behavior,
  });
}

export function scrollContainerToSectionBottom(
  container: HTMLElement,
  section: HTMLElement,
  behavior: ScrollBehavior = "smooth",
  excludeSelector?: string,
) {
  const excluded = excludeSelector
    ? section.querySelector<HTMLElement>(excludeSelector)
    : null;
  const excludedHeight = excluded?.offsetHeight ?? 0;
  const contentBottom =
    section.offsetTop + section.offsetHeight - excludedHeight;
  const targetScroll = Math.max(0, contentBottom - container.clientHeight);

  container.scrollTo({ top: targetScroll, behavior });
}
