```bash
# Install Poetry
curl -sSL https://install.python-poetry.org | python3 -

# Clone repo
git clone https://github.com/igor-policee/kube_scan.git
cd ./kube_scan
chmod u+x ./kube_scan.py

# Install dependencies
${HOME}/.local/bin/poetry install

# Activate virtual environment
${HOME}/.local/bin/poetry shell

# Run script
python ./kube_scan.py -h
```