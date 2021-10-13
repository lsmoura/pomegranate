document.addEventListener("DOMContentLoaded", main);

/**
 * @typedef Movie
 * @property {Array<string>} genres
 * @property {{Posters: Array<String>}} images
 * @property {string} imdb_id
 * @property {string} released release date in the format YYYY-MM-DD
 * @property {number} runtime
 * @property {Array<string>} titles
 * @property {number} tmdb_id
 * @property {number} year
 */

/**
 * @typedef MovieAddResponse
 * @property {string} message
 * @property {string} title
 * @property {string} overview
 */

/**
 * @typedef MovieResponseItem
 * @property {string} imdb_id
 * @property {string} title
 *
 * @typedef {Array<MovieResponseItem>} MovieListResponse
 */

const movieTitleID = "movie-title";
const movieSearchResultsID = "movie-search-results";
const clearSearchButtonID = "clear-search"
const moviesListID = "movies-list"

let debouncedSearch = null;

function updateMovieSearch(event) {
  const str = event.target.value;

  clearTimeout(debouncedSearch)
  debouncedSearch = setTimeout(function() {
    executeMovieSearch(str);
  }, 500);
}

/**
 * populateMovieResult Creates 'li' items from the movie array parameter and injects on the dst html element
 * @param dst {HTMLElement}
 * @param movies {Array<Movie>}
 */
function populateMovieResult(dst, movies) {
  if (movies === null || (Array.isArray(movies) && movies.length <= 0)) {
    dst.innerHTML = "no results";
    return;
  }

  movies.forEach(function (movie) {
    const movieEntry = document.createElement("li");
    movieEntry.classList.add("movie-search-result")

    const movieTitle = document.createElement("span");
    movieTitle.classList.add("title")
    movieTitle.innerHTML = movie.titles[0]
    movieEntry.appendChild(movieTitle)

    const movieImdbID = document.createElement("span");
    movieImdbID.classList.add("imdb-id")
    movieImdbID.innerHTML = movie.imdb_id
    movieEntry.appendChild(movieImdbID)

    const movieReleaseDate = document.createElement("span");
    movieReleaseDate.classList.add("release-date")
    movieReleaseDate.innerHTML = movie.released
    movieEntry.appendChild(movieReleaseDate)

    if (movie.imdb_id) {
      const movieAddButton = document.createElement("button")
      movieAddButton.innerHTML = "add"
      movieAddButton.addEventListener("click", function(event) { handleMovieAdd(movie.imdb_id) })
      movieEntry.appendChild(movieAddButton)
    }

    dst.appendChild(movieEntry);
  });
}

/**
 * Execute API query to search for movies
 * @param str {string}
 * @returns {Promise<Array<Movie>>}
 */
function apiMovieSearch(str) {
  return fetch("/movie/search?q=" + str)
    .then(function (result) { return result.json(); })
    .then(function (result) { return result.movies; })
}

/**
 * Adds a movie to the database
 * @param imdbId {string}
 * @returns {Promise<MovieAddResponse>}
 */
function apiMovieAdd(imdbId) {
  return fetch("/movie/add?identifier=" + imdbId)
    .then(function (result) { return result.json(); });
}

/**
 * Returns the list of movies that are being managed
 * @returns {Promise<MovieListResponse>}
 */
function apiMovieList() {
  return fetch("/movie/list")
    .then(function (result) { return result.json(); });
}

let latestMovieSearch = "";
function executeMovieSearch(str) {
  latestMovieSearch = str;
  const element = document.getElementById(movieSearchResultsID);
  if (!element) {
    console.error("search results container not found: '" + movieSearchResultsID + "'");
    return;
  }
  if (str === "") {
    element.innerHTML = "";
    return;
  }

  element.innerHTML = "Loading...";

  apiMovieSearch(str)
    .then(function (movies) {
      if (str !== latestMovieSearch) return;
      const movieList = document.createElement("ul");
      populateMovieResult(movieList, movies);
      element.innerHTML = "";
      element.appendChild(movieList);
    })
    .catch(function(err) {
      console.error(err);
    })
}

function clearSearch() {
  executeMovieSearch("");

  const movieTitleInputElement = document.getElementById("movie-title");
  if (!movieTitleInputElement) {
    console.error("cannot find the movie title input");
  } else {
    movieTitleInputElement.value = "";
  }
}

function handleMovieAdd(imdbId) {
  apiMovieAdd(imdbId)
    .then(function(response) {
      console.log(response)

      return updateMovieList();
    });
}

/**
 * Helper function to attach event listeners to elements
 * @param elementId {string}
 * @param type {string}
 * @param callback {(event: unknown) => void}
 */
function registerHandler(elementId, type, callback) {
  const element = document.getElementById(elementId);
  if (!element) {
    console.error("cannot find element with id '" + elementId + "'");
    return;
  }

  element.addEventListener(type, callback);
}

function updateMovieList() {
  const element = document.getElementById(moviesListID)
  if (!element) {
    console.error("cannot find element with id '" + moviesListID + "'");
    return;
  }

  element.innerHTML = "Loading...";
  apiMovieList()
    .then(function (list) {
      element.innerHTML = "";
      if (list) {
        list.forEach(function (item) {
          const entry = document.createElement("div");
          entry.innerHTML = item.title + " (" + item.imdb_id + ")";
          element.appendChild(entry)
        });
      }
    });
}

function main() {
  console.log("Pomegranate initializing...");
  registerHandler(movieTitleID, "keyup", updateMovieSearch);
  registerHandler(clearSearchButtonID, "click", clearSearch);

  updateMovieList();
}

