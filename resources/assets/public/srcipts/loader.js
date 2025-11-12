var EDUCATION_BLOCK = document.getElementById("education-block");
var MISSION_BLOCK = document.getElementById("mission-block");
var JOB_BLOCK = document.getElementById("job-block");
var CHARITY_BLOCK = document.getElementById("charity-block");

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

function load_mission(json_obj) {
    let fragment = document.createDocumentFragment();
    let p = document.createElement("p");
    p.textContent = json_obj.text;
    fragment.appendChild(p);

    MISSION_BLOCK.appendChild(fragment);
}

function load_education(json_obj) {
    let fragment = document.createDocumentFragment();
    json_obj.forEach(school_entry => {
        let line = `${school_entry.degree} | ${school_entry.attended}`;
        make_content(fragment, school_entry.name, school_entry.link, school_entry.image, [line]);
    })

    EDUCATION_BLOCK.appendChild(fragment);
}

function load_jobs(json_obj) {
    let fragment = document.createDocumentFragment();
    json_obj.forEach(job_entry => {
        make_content(fragment, job_entry.company, job_entry.link, job_entry.image, [job_entry.role, job_entry.attended, job_entry.description, job_entry.experience_gained]);
    })

    JOB_BLOCK.appendChild(fragment);
}

function load_charities(json_obj) {
    let fragment = document.createDocumentFragment();
    json_obj.forEach(charity_entry => {
        make_content(fragment, charity_entry.name, charity_entry.link, charity_entry.image, [charity_entry.description, charity_entry.contribution]);
    });

    CHARITY_BLOCK.appendChild(fragment);
}

async function cache_or_load(cache_name, json_file, load_function) {
    const cached = localStorage.getItem(cache_name);
    if (cached) {
        load_function(JSON.parse(cached));
        // force update (later This should happen less often maybe only every 30 minutes or so, this could also be done lazily by just free-ing the cache)
        const response = await fetch(json_file, { cache: "force-cache" });
        const data = await response.json();
        localStorage.setItem(cache_name, JSON.stringify(data));
        return;
    }

    try {
        const response = await fetch(json_file, { cache: "force-cache" });
        const data = await response.json();
        load_function(data);
        localStorage.setItem(cache_name, JSON.stringify(data));
    } catch (error) {
        console.error("Failed to load website data: ", error);
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    cache_or_load("mission-data", "data/mission.json", load_mission);
    cache_or_load("education-data", "data/schools.json", load_education);
    cache_or_load("job-data", "data/jobs.json", load_jobs);
    cache_or_load("charity-data", "data/charities.json", load_charities);
});