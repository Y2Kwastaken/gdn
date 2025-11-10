var EDUCATION = document.getElementById("education-block");

/*
*         <div class="content-general">
            <div class="content-general-header">
                <a href="https://google.com" target="_blank" rel="noopener noreferrer"> <!-- Change cursor also -->
                    <img src="imgs/college.png" />
                </a>
                <h3>School Institution 1</h3>
            </div>
            <p>Bachelors of Science, Computer Science | 20XX-20XX</p>
        </div>
*/

function make_content(append_to, title_string, link_string, image, description_string_iterable) {
    let parent_node = document.createElement("div");
    parent_node.className = "content-general";
    parent_node.innerHTML = `
    <div class="content-general-header">
        <a href="${link_string}" target="_blank" rel="noopener noreferrer"> <!-- Change cursor also -->
            <img loading="lazy" src="${image}" />
        </a>
        <h3>${title_string}</h3>
    </div>
    `

    description_string_iterable.forEach(element => {
        let tmp = document.createElement("p");
        tmp.textContent = element;
        parent_node.appendChild(tmp);
    });

    append_to.appendChild(parent_node);
}

function load_education(json_obj) {
    let fragment = document.createDocumentFragment();
    json_obj.forEach(school_entry => {
        let line = `${school_entry.degree} | ${school_entry.attended}`;
        make_content(fragment, school_entry.name, school_entry.link, school_entry.image, [line]);
    })

    EDUCATION.appendChild(fragment);
}

document.addEventListener("DOMContentLoaded", async () => {
    const cached = localStorage.getItem("education-data");
    if (cached) {
        load_education(JSON.parse(cached));
        return;
    }

    // The preload ensures this returns nearly instantly
    try {
        const response = await fetch("data/schools.json", { cache: "force-cache" });
        const data = await response.json();
        load_education(data);
        localStorage.setItem("education-data", JSON.stringify(data));
    } catch (error) {
        console.error("Education Load Error:", error);
    }
});