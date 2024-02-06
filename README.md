# Talisman GitHub Pages

The pages are written using just-the-docs template in Jekyll. GitHub Pages
automatically builds and serves the site. To best match the deployed
environment, use something like [rbenv](https://github.com/rbenv/rbenv) to match
the Ruby version that [GitHub Pages uses](https://pages.github.com/versions/).

To locally run the jekyll server:

* `bundle config set path vendor/bundle`
* `bundle install`
* `bundle exec jekyll serve`

Your pages will start on `localhost:4000/talisman`
