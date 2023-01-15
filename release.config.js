module.exports = {
  "branches":  [
    { "name": "main" },
    { "name": "beta", "channel": "beta", "prerelease": "beta" },
  ],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    ["@semantic-release/exec", {
      "publishCmd": "make VERSION=${nextRelease.version} release -j 3"
    }],
    ["@semantic-release/github", {
      "assets": [
        "sablier*"
      ]
    }],
    ["@semantic-release/exec", {
      "publishCmd": "make LAST=${lastRelease.version} NEXT=${nextRelease.version} update-doc-version update-doc-version-middleware"
    }],
    ["@semantic-release/git", {
      "assets": ["**/*.{md,yml}"],
      "message": "docs(release): update doc version from ${lastRelease.version} to ${nextRelease.version} [skip ci]"
    }]
  ]
}