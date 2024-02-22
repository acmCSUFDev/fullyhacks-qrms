export function updateSearch(search) {
  const targets = document.querySelectorAll(search.dataset.search);
  const selectors = search.dataset.searchSelectors
    .split(",")
    .map((e) => e.trim());
  if (!selectors) {
    throw new Error(`Missing data-search-selectors attribute on search input`);
  }

  const query = search.value.toLowerCase();

  targets.forEach((e) => {
    const str = selectors.reduce((acc, selector) => {
      const value = e.querySelector(selector).textContent.toLowerCase();
      console.log(selector, value, acc + value);
      return acc + value;
    });

    if (str.includes(query)) {
      e.style.display = "block";
    } else {
      e.style.display = "none";
    }
  });
}

document
  .querySelectorAll('input[type="search"][data-search]')
  .forEach((search) => {
    search.disabled = false;
    search.addEventListener("input", () => updateSearch(search));
  });
