#!/bin/bash

cd ../kubectl-ice
git checkout pages

cd ../ice-web

JEKYLL_ENV=production bundle exec jekyll build -d ../kubectl-ice/docs
