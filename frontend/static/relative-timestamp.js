const Duration = {
  second: 1000,
  minute: 60 * 1000,
  hour: 60 * 60 * 1000,
  day: 24 * 60 * 60 * 1000,
};

const createRelativeTimeFormatter = (style) =>
  new Intl.RelativeTimeFormat(undefined, {
    style,
    numeric: "auto",
  });

const relativeTimes = {
  short: createRelativeTimeFormatter("short"),
  long: createRelativeTimeFormatter("long"),
};

document.querySelectorAll("time.relative").forEach((time) => {
  let formatter = relativeTimes.long;
  if (time.classList.contains("short")) {
    formatter = relativeTimes.short;
  }

  const date = new Date(time.getAttribute("datetime"));

  let delta = date.getTime() - Date.now();
  let unit;
  if (Math.abs(delta) < Duration.minute) {
    delta = Math.round(delta / Duration.second);
    unit = "second";
  } else if (Math.abs(delta) < Duration.hour) {
    delta = Math.round(delta / Duration.minute);
    unit = "minute";
  } else if (Math.abs(delta) < Duration.day) {
    delta = Math.round(delta / Duration.hour);
    unit = "hour";
  } else {
    delta = Math.round(delta / Duration.day);
    unit = "day";
  }

  time.textContent = formatter.format(delta, unit);
});

document.querySelectorAll("time.localize").forEach((time) => {
  const date = new Date(time.getAttribute("datetime"));
  const localized = date.toLocaleString(undefined, {
    dateStyle: "full",
    timeStyle: "long",
  });
  time.title = localized;
  time.textContent = localized;
});
