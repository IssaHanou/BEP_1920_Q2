import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="scclib",
    version="0.0.6",
    author="Raccoon Serious Games",
    author_email="BEP@raccoon.games",
    description="Client library for S.C.I.L.E.R.",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/IssaHanou/BEP_1920_Q2",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: GNU General Public License v3 (GPLv3)",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.5",
)
