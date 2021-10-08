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

const movieTitleID = "movie-title";
const movieSearchResultsID = "movie-search-results";

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
    const movieEntry = document.createElement('li');

    const movieTitle = document.createElement('span');
    movieTitle.innerHTML = movie.titles[0]
    movieEntry.appendChild(movieTitle)

    const movieImdbID = document.createElement('span');
    movieImdbID.innerHTML = movie.imdb_id
    movieEntry.appendChild(movieImdbID)

    const movieReleaseDate = document.createElement('span');
    movieReleaseDate.innerHTML = movie.released
    movieEntry.appendChild(movieReleaseDate)

    dst.appendChild(movieEntry);
  });
}

let latestMovieSearch = "";
function executeMovieSearch(str) {
  latestMovieSearch = str;
  const element = document.getElementById(movieSearchResultsID);
  if (str === "") {
    element.innerHTML = "";
    return;
  }

  element.innerHTML = "Loading...";


  fetch("/movie/search?q=" + str)
    .then(function (result) { return result.json() })
    .then(function (result) {
      if (str !== latestMovieSearch) return;
      const movieList = document.createElement('ul');
      populateMovieResult(movieList, result.movies)
      element.innerHTML = ""
      element.appendChild(movieList)
    })
    .catch(function(err) {
      console.error(err);
    })
}

function main() {
  console.log("Pomegranate initializing...");
  movieSearchElement = document.getElementById(movieTitleID);
  if (movieSearchElement) {
    movieSearchElement.addEventListener("keyup", updateMovieSearch)
  } else {
    console.error("cannot find element with id '" + movieTitleID + "'")
  }
}

