direnv_store_dir=${DIRENV_STORE:-$direnv_config_dir/store}
mkdir -p $direnv_store_dir

# store current source_env as _source_env
eval "`declare -f source_env | sed '1s/.*/_&/'`"

source_env() {
  _source_env $@

  # create a symlink to the direnv layout dir inside the store
  direnv_real_layout_dir=$(realpath $(direnv_layout_dir))
  if [[ -d "$direnv_real_layout_dir" ]]; then
    direnv_layout_dir_hash=$(echo -n $direnv_real_layout_dir | sha256sum | cut -d ' ' -f 1)
    direnv_store_path=$direnv_store_dir/$direnv_layout_dir_hash
    
    echo "direnv: storing $direnv_real_layout_dir -> $direnv_store_path"
    
    ln -sf $direnv_real_layout_dir $direnv_store_path
    touch -mh $direnv_store_path
  fi
}