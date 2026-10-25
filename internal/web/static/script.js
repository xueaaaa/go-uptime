const labels = {available: "Up", unavailable: "Down", unknown: "Unknown"};
const dots = {available: "up", unavailable: "down", unknown: "unknown"};

const esc = s => String(s).replace(/[&<>"']/g, c => (
    {"&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;"}[c]
));

const dot = k => `<div class="status-indicator status-indicator--${dots[k]}"></div>`;

const tbody = document.querySelector("#urls-table-body");
const search = document.querySelector("#url-search");

function applyFilter() {
    const q = search.value.trim().toLowerCase();

    tbody.querySelector(".no-results")?.remove();

    const rows = [...tbody.querySelectorAll("tr")];
    let visible = 0;

    for (const row of rows) {
        const url = row.cells[0].textContent.toLowerCase();
        const match = url.includes(q);
        row.style.display = match ? "" : "none";
        if (match) visible++;
    }

    if (rows.length > 0 && visible === 0) {
        tbody.insertAdjacentHTML("beforeend",
            `<tr class="no-results"><td colspan="4">Nothing found</td></tr>`);
    }
}

search.addEventListener("input", applyFilter);

new EventSource("/events/sites").onmessage = e => {
    const d = JSON.parse(e.data);
    if (!Array.isArray(d.sites)) return;

    tbody.innerHTML = d.sites.map(s => `
        <tr>
            <td><a href="/site/${esc(s.id)}">${esc(s.url)}</a></td>
            <td>
                <div class="status-item">
                    ${dot(s.class)}
                    <span>${labels[s.class] ?? "Unknown"}</span>
                </div>
            </td>
            <td>${esc(s.interval)}</td>
            <td>${esc(s.lastCheck)}</td>
        </tr>
    `).join("");

    const total = d.total === 0 ? "No sites" : d.total === 1 ? "1 site" : `${d.total} sites`;
    const item = (k, n) => `<span>|</span><div class="summary-item">${dot(k)}<span>${n}</span></div>`;

    document.querySelector(".summary-info").innerHTML =
        `<div class="summary-item"><span class="summary-content">${total}</span></div>` +
        (d.total !== 0
            ? item("available", d.available) + item("unavailable", d.unavailable) + item("unknown", d.unknown)
            : "");

    applyFilter();
};
