// load categories on page load
$.get("http://localhost:8080/feed/categories", function (data) {
	data.forEach((v, i) => {
		$(".categories .list").append(`<a class="waves-effect waves-light btn-small orange ${i == 0 ? "disabled" : ""}" data-category="${v}">${v}</a>`)
	})

	fetchRSSs(data[0]);
});

// change categories
$(".categories .list").on("click", "a", (e) => {
	$(".categories .list .disabled").removeClass("disabled");
	$(e.currentTarget).addClass("disabled");
	fetchRSSs($(e.currentTarget).attr("data-category"));
});

// delete rss
$(".rsss").on("click", ".delete", (e) => {
	const currentTarget = $(e.currentTarget);
	const thisCard = currentTarget.closest(".card");
	const url = encodeURIComponent(thisCard.find("a").attr("href"));

	const settings = {
		"async": true,
		"crossDomain": true,
		"url": `http://localhost:8080/rss/${url}`,
		"method": "DELETE",
	}

	$.ajax(settings).done(() => thisCard.fadeOut(200)).fail((e) => console.log(e));
});

// fetchRSSs gets a list of rss, hides the old and creates the new one
function fetchRSSs(category) {
	$.get(`http://localhost:8080/rss/category/${category}`, function (data) {
		$("main .rsss").html("");

		if (!data.rsss) return;

		data.rsss.forEach(v => {
			const subtitle = v.subtitle  ? `<p>${v.subtitle}</p>` : '';

			$("main .rsss").append(
				`<div class="card">
					<a href="${v.url}" target="_blank">
				      <div class="card-content white-text text-darken-1">
				        <span class="card-title">${v.title}</span>
				        
				         ${subtitle}
				      </div>
				    </a>
				    
				    <div class="card-action yellow-text text-darken-3">
				      <span>${v.source.replace('https://', '').replace('http://', '')}</span>
				      
				      <i class="delete tiny material-icons">delete</i>
				    </div>
				</div>`
			);
		});
	});
}
