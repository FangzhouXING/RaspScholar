const photoHolder = document.getElementById("photoHolder");
const name = document.getElementById("name");
const university = document.getElementById("university");
const area = document.getElementById("area");
const hIndex = document.getElementById("hIndex");

let slideIndex = 0;

carousel();
updateFromServer();

function updateFromServer() {
	updateBio(fetchBioFromServer());
	// updateMostCitedPapers(fectchRecentPublications());
	updateRecentPublications(fectchRecentPublications());

	// 2 seconds
	setTimeout(carousel, 2000); 
}

function fetchHelper(path) {
	fetch(path, {
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: "GET",
    }).then((response) => {
        response.text().then(function (data) {
            // console.log(">>>|"+data+"|<<<")
            let result = JSON.parse(data);
            console.log(result)
            return result;
        });
    });
}

function fetchBioFromServer() {
	return fetchHelper("/bascinfo/");
}

function updateBio(basicInfo) {
	// update photo
	// TODO
	// update name
	name.innerText = basicInfo["Name"];
	// update university
	university.innerText = basicInfo["RankAndSchool"];
	// update focus area
	university.innerText = basicInfo["Focus"].join(", ");
	// update h-index
	// TODO
}


function fetchMostCitedPublications() {

}

function updateMostCitedPapers() {

}

function fectchRecentPublications() {
	return fetchHelper("/recentpapers/");
}

function updateRecentPublications(papers) {
	let recentPublicationsList = document.querySelectorAll("#recentPublications > slideListItem");
	let i = 0;
	for (i = 0; i < recentPublicationsList.length; i++) {
		let title = "";
		if (papers["Papers"].length < i) {
			title = papers["Papers"][i]["Title"];
		}
		recentPublicationsList[i].innerText = title;
	}
}

function carousel() {
	var i;
	var x = document.getElementsByClassName("slide");
	for (i = 0; i < x.length; i++) {
		x[i].style.display = "none"; 
	}
	slideIndex++;
	if (slideIndex > x.length) {
		slideIndex = 1;
	} 
	x[slideIndex-1].style.display = "block"; 
	// setTimeout(carousel, 1000); 
}

