from setuptools import setup, find_packages
setup(
    name = "Whistle-log",
    version = "0.1",
    packages = find_packages(),
    scripts = ['web_server.py'],

    # Project uses reStructuredText, so ensure that the docutils get
    # installed or upgraded on the target machine
    install_requires = ['flask', 'fuzzywuzzy', 'pymongo', 'pytz', 'google-api-python-client', 'requests'],

    package_data = {
        # If any package contains *.txt or *.rst or .md files, include them:
        '': ['*.txt', '*.rst','*.md']
    },

    # metadata for upload to PyPI
    author = "Roopak Parikh",
    author_email = "rparikh@platform9.com",
    description = "Whistle is a log aggregating solution for Platform9 alerts",
    license = "PSF",
    keywords = "",
    url = "http://whistle.platform9.horse/",   # project home page, if any
    # could also include long_description, download_url, classifiers, etc.
)
