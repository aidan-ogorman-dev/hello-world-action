# Label checker docker action

This action checks added and modified files for k8s manifests, and updates them with required labels.

## Inputs

None

## Outputs

Modified files are written to ${{ github.workspace }}.

## Example usage

uses: actions/hello-world-docker-action@v2
with:
  who-to-greet: 'Mona the Octocat'
