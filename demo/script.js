async function main() {
	const go = new Go();

	console.log("Loading WebAssembly module...");
	const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
	console.log("WebAssembly module loaded.");
	go.run(result.instance);

	document.getElementById("src").addEventListener("input", debounce(function() {
		const src = document.getElementById("src").value;
		const t0 = performance.now();
		const html = window.xatenaToHTML(src);
		const t1 = performance.now();
		console.log(`xatenaToHTML took ${t1 - t0} ms`);
		const resultDiv = document.getElementById("result");
		console.log(html)
		resultDiv.innerHTML = html;
		setDataTagRecursive(resultDiv);
	}, 500));

	// 初期状態を反映
	document.getElementById("src").dispatchEvent(new Event("input"));
}

// ディレイ付き実行（debounce）関数
function debounce(fn, delay) {
	let timer = null;
	return function(...args) {
		clearTimeout(timer);
		timer = setTimeout(() => fn.apply(this, args), delay);
	};
}

function setDataTagRecursive(el) {
	if (el.nodeType === Node.ELEMENT_NODE) {
		el.setAttribute('data-tag', el.tagName.toLowerCase());
		for (const child of el.children) {
			setDataTagRecursive(child);
		}
	}
}

main();
