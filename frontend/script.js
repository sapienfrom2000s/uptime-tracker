const timeline = document.getElementById("timeline");

// Fetch data from your Go API
fetch("http://localhost:8080/v1/trackers/status")
  .then((response) => response.json())
  .then((data) => {
    timeline.innerHTML = ""; // clear old elements

    data.forEach((up) => {
      const rect = document.createElement("div");
      rect.classList.add("rect");

      // map boolean -> color
      if (up === true) rect.classList.add("green");
      else if (up === false) rect.classList.add("red");
      else rect.classList.add("white");

      timeline.appendChild(rect);
    });
  })
  .catch((err) => {
    console.error("Error fetching statuses:", err);
  });
