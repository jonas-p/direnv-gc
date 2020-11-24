direnv_store_dir=${DIRENV_STORE:-$direnv_config_dir/store}
mkdir -p $direnv_store_dir

# store current source_env as _source_env
eval "$(declare -f source_env | sed '1s/.*/_&/')"

# store current direnv_layout_path as _direnv_layout_path
eval "$(declare -f direnv_layout_dir | sed '1s/.*/_&/')"

direnv_gc_temp_file=$(mktemp)

# direnv_layout_dir could be called multiple times when running
# direnv allow and can also be called from subshells etc so we use
# a temp file to store direnv locations we can update the links later
#
# it's important that this happens last in the direnvrc because the
# user might override it
direnv_layout_dir() {
	local dir=$(_direnv_layout_dir $@)
	echo $dir >>$direnv_gc_temp_file
	echo $dir
}

direnv_gc_link() {
	# create a symlink to the direnv layout dir inside the store
	direnv_real_layout_dir=$(realpath $@)
	if [[ -d "$direnv_real_layout_dir" ]]; then
		direnv_layout_dir_hash=$(echo -n $direnv_real_layout_dir | sha256sum | cut -d ' ' -f 1)
		direnv_store_path=$direnv_store_dir/$direnv_layout_dir_hash

		echo "direnv: storing $direnv_real_layout_dir -> $direnv_store_path"

		ln -sf $direnv_real_layout_dir $direnv_store_path
		touch -mh $direnv_store_path
	fi
}

source_env() {
	# restore old source env, this is function to be called once, and
	# not from any child/parent direnvs (i.e. using source_up)
	eval "$(declare -f _source_env | sed '1s/_source_env/source_env/')"
	source_env $@

	# link unique direnvs to the store
	sort -u $direnv_gc_temp_file | while read line; do
		direnv_gc_link $line
	done
}
