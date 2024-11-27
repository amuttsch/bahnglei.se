# bahnglei.se

## Importer

Create the importer machine with a volume:

```sh
fly vol create -s 10 --region ams osm_data
```

```sh
fly m create --schedule monthly --restart no --vm-cpus 2 --vm-memory 1024 -r ams -v osm_data:/osm -n importer --entrypoint "/app/bahngleise import" $(fly app releases --image -a bahngleise | grep complete | awk 'BEGIN {FS="\t"}; {print $6}' | head -1)
```

Run it with:

```sh
fly m start $(fly m list | grep importer | awk -F' ' '{print $1}')
```

Update it with:

```sh
fly m update $(fly m list | grep importer | awk -F' ' '{print $1}') --image $(fly app releases --image -a bahngleise | grep complete | awk 'BEGIN {FS="\t"}; {print $6}' | head -1)
```
