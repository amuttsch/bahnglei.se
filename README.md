# bahnglei.se

## Importer

Create the importer machine:

```sh
fly m create --schedule monthly --vm-cpus 4 --vm-memory 2048 -n importer --entrypoint "/app/bahngleise import" $(fly app releases --image -a bahngleise | grep complete | awk 'BEGIN {FS="\t"}; {print $6}' | head -1)
```

Run it with:

```sh
fly m start $(fly m list | grep importer | awk -F' ' '{print $1}')
```

Update it with:

```sh
fly m update $(fly m list | grep importer | awk -F' ' '{print $1}') --image $(fly app releases --image -a bahngleise | grep complete | awk 'BEGIN {FS="\t"}; {print $6}' | head -1)
```
