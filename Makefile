
# 从 svn 拉取文件 (map, excel)
svn_checkout:
	@ echo "cmd: [svn-checkout] staring... [$(PWD)]"
	@ if [ -f svn.env ]; then source svn.env; fi; \
	./script/svn_checkout.bash $$SVN_USER $$SVN_PASS  $(SVN_REPO_MAP) $(TEMP_DIR_MAP)
	@ mkdir -p $(DIR_MAP)
	@ rm -rf $(DIR_MAP)/*server.json
	@ cp $(TEMP_DIR_MAP)/*server.json $(DIR_MAP)
	@ if [ -f svn.env ]; then source svn.env; fi; \
	./script/svn_checkout.bash $$SVN_USER $$SVN_PASS  $(SVN_REPO_EXCEL) $(TEMP_DIR_EXCEL)
	@ mkdir -p $(DIR_EXCEL)
	@ rm -rf $(DIR_EXCEL)/*.xlsx
	@ find "$(TEMP_DIR_EXCEL)" -type f -name "*.xlsx" ! -name "I18n*" -exec cp {} "$(DIR_EXCEL)" \;
	@ echo "cmd: [svn-checkout] done.\n"
