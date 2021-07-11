const photoHolder = document.getElementById("photoHolder");
const name = document.getElementById("name");
const university = document.getElementById("university");
const area = document.getElementById("area");
const hIndex = document.getElementById("hIndex");

let slideIndex = 0;

showSlides();
updateFromServer();

function updateFromServer() {
	updateBio(fetchBioFromServer());
	// updateMostCitedPapers(fectchRecentPublications());
	updateRecentPublications(fectchRecentPublications());

	// 60 seconds
	setTimeout(updateFromServer, 60000); 
}

function fetchHelper(path) {
	return fetch(path, {
		headers: {
			'Accept': 'application/json',
			'Content-Type': 'application/json'
		},
		method: "GET",
	}).then((response) => {
		return response.json();
	}
	);
}

function fetchBioFromServer() {
	return fetchHelper("/bascinfo/");
}

function updateBio(basicInfoPromise) {
	basicInfoPromise.then(
		(basicInfo) => {
			// console.log("basicInfo");
			// console.log(basicInfo);
			// console.log(basicInfo["Name"]);
			// update photo
			// TODO
			// update name
			// console.log(basicInfo);
			name.innerText = basicInfo["Name"];
			// update university
			// university.innerText = basicInfo["RankAndSchool"];
			// update focus area
			area.innerText = basicInfo["Focus"].join(", ");
			// update h-index
			hIndex.innerText = basicInfo["HIndex"];
			// TODO
		});
	
}


function fetchMostCitedPublications() {

}

function updateMostCitedPapers() {

}

function fectchRecentPublications() {
	return fetchHelper("/recentpapers/");
}

function updateRecentPublications(papersPromise) {
	papersPromise.then(
		papers => {
			let recentPublicationsList = document.querySelectorAll("#recentPublications .slideListItem");
			let i = 0;
			// console.log(recentPublicationsList.length);
			// console.log("recentPublicationsList");
			// console.log(papers);
			for (i = 0; i < recentPublicationsList.length; i++) {
				// console.log(i);
				let title = "";
				if (i < papers["Papers"].length) {
					title = papers["Papers"][i]["Title"];
					// console.log(title);
				}
				recentPublicationsList[i].innerText = title;
			}
		}
		);
	
}

function showSlides() {
  let i;
  let slides = document.getElementsByClassName("slide");
  let dots = document.getElementsByClassName("dot");
  for (i = 0; i < slides.length; i++) {
    slides[i].style.display = "none";  
  }
  slideIndex++;
  if (slideIndex > slides.length) {slideIndex = 1}    
  for (i = 0; i < dots.length; i++) {
    dots[i].className = dots[i].className.replace(" active", "");
  }
  slides[slideIndex-1].style.display = "block";  
  dots[slideIndex-1].className += " active";
  setTimeout(showSlides, 10000); // Change image every 10 seconds
}

