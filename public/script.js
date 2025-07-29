const qs = (selector) => document.querySelector(selector);

function main() {
  const urlParams = new URLSearchParams(location.search);
  const redirectURI = urlParams.get("redirect_to") ?? "https://monitor.deps.me";
  qs("#submit").onclick = async () => {
    const password = qs("#password-input").value;
    const res = await fetch("/login", {
      method: "POST",
      body: password,
    });
    if (res.status == 200) {
      location.href = redirectURI;
    } else {
      console.log("retry");
    }
  }
}

main()
