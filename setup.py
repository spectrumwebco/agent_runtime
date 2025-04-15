from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

with open("requirements.txt", "r", encoding="utf-8") as fh:
    requirements = fh.read().splitlines()

setup(
    name="fine-tune",
    version="0.1.0",
    author="Spectrum Web Co",
    author_email="oveshen.govender@gmail.com",
    description="Web scraper for collecting solved issues from GitHub and Gitee repositories to fine-tune Llama 4 models",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/spectrumwebco/Fine-Tune",
    packages=find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.8",
    install_requires=requirements,
)
