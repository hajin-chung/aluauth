const qs = (selector) => document.querySelector(selector);

async function handleSubmit() {
  const urlParams = new URLSearchParams(location.search);
  const redirectURI = urlParams.get("redirect_to") ?? "https://monitor.deps.me";

  const password = qs("#password-input").value;
  const res = await fetch("/login", {
    method: "POST",
    body: password,
  });
  if (res.status == 200) {
    location.href = redirectURI;
  } else {
    qs("#password-input").classList.add("wrong");
  }
}

function main() {
  qs("#password-input").addEventListener("keydown", (evt) => {
    if (evt.key === "Enter") handleSubmit();
  });
  qs("#submit").onclick = handleSubmit;
}

main()
